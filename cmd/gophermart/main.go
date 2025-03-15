package main

import (
	"context"
	"database/sql"
	"errors"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/spf13/pflag"

	"github.com/vilasle/gophermart/internal/logger"
	_middleware "github.com/vilasle/gophermart/internal/middleware"
	"github.com/vilasle/gophermart/internal/service/gophermart/accrual"
	"github.com/vilasle/gophermart/internal/service/gophermart/authorization"
	"github.com/vilasle/gophermart/internal/service/gophermart/order"
	"github.com/vilasle/gophermart/internal/service/gophermart/withdrawal"

	httpRep "github.com/vilasle/gophermart/internal/repository/gophermart/http"
	pgRep "github.com/vilasle/gophermart/internal/repository/gophermart/postgresql"

	"github.com/vilasle/gophermart/internal/controller/gophermart"
)

type cliArgs struct {
	addr        string
	accrualAddr string
	dbURI       string
	debug       bool
}

func initCli() cliArgs {
	args := cliArgs{}
	pflag.StringVarP(&args.addr, "address", "a", ":8080", "address to listen on")
	pflag.StringVarP(&args.dbURI, "database", "d", "", "database url e.g postgres://postgres:postgres@localhost:5432/postgres")
	pflag.StringVarP(&args.accrualAddr, "accrual", "r", "", "accrual endpoint")

	pflag.BoolVarP(&args.debug, "debug", "D", false, "enable debug message")
	pflag.Parse()

	args.addr = getEnv("RUN_ADDRESS", args.addr)
	args.dbURI = getEnv("DATABASE_URI", args.dbURI)
	args.accrualAddr = getEnv("ACCRUAL_SYSTEM_ADDRESS", args.accrualAddr)

	return args
}

func getEnv(key, fallback string) string {
	result := fallback
	if value, _ := os.LookupEnv(key); value != "" {
		result = value
	}
	return result
}

func main() {
	args := initCli()

	initLogger(args)

	if err := checkArgs(args); err != nil {
		logger.Error("invalid arguments", "error", err)
		pflag.Usage()
		os.Exit(1)
	}

	accrualURL, err := url.Parse(args.accrualAddr)
	if err != nil {
		logger.Error("can not parse accrual endpoint address", "error", err)
		os.Exit(1)
	}

	db, err := sql.Open("pgx", args.dbURI)
	if err != nil {
		logger.Error("connecting to database failed", "url", args.dbURI, "error", err)
		os.Exit(1)
	}
	defer db.Close()

	dbRep, err := pgRep.NewPostgresqlGophermartRepository(db)
	if err != nil {
		logger.Error("can not init gophermart repository", "error", err)
		os.Exit(1)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	orderSvc := createOrderService(dbRep, accrualURL)
	orderSvc.Start(ctx)

	ctrl := newController(dbRep, orderSvc)

	mux := newMux(ctrl)

	server := newServer(mux, args.addr)
	defer server.Close()

	s := signalSubscription()
	logger.Info("run server", "addr", args.addr)
	go run(server, s)

	<-s

	shutdown(ctx, server)

	cancel()

	time.Sleep(time.Millisecond * 2500)
}

func initLogger(args cliArgs) {
	if args.debug {
		logger.Init(os.Stdout, logger.DebugLevel)
	} else {
		logger.Init(os.Stdout, logger.InfoLevel)
	}
}

func checkArgs(args cliArgs) error {
	errs := make([]error, 3)
	if args.dbURI == "" {
		errs = append(errs, errors.New("database url is required"))
	}

	if args.addr == "" {
		errs = append(errs, errors.New("address is required"))
	}

	if args.accrualAddr == "" {
		errs = append(errs, errors.New("accrual endpoint is required"))
	}

	return errors.Join(errs...)
}

func createOrderService(pgRepository pgRep.PostgresqlGophermartRepository, accrualURL *url.URL) order.OrderService {
	accrualSvc := accrual.NewAccrualService(
		httpRep.NewAccrualRepository(accrualURL),
	)

	orderSvc := order.NewOrderService(order.OrderServiceConfig{
		OrderRepository:        pgRepository,
		AccrualService:         accrualSvc,
		WithdrawalRepository:   pgRepository,
		RetryOnError:           time.Second * 10,
		AttemptsGettingAccrual: 3,
	})

	return orderSvc
}

func newController(pgRepository pgRep.PostgresqlGophermartRepository, svc order.OrderService) gophermart.Controller {
	withdrawalSvc := withdrawal.NewWithdrawalService(pgRepository)

	authSvc := authorization.NewAuthorizationService(pgRepository)

	return gophermart.Controller{
		AuthSvc:     authSvc,
		OrderSvc:    svc,
		WithdrawSvc: withdrawalSvc,
	}
}

func newMux(ctrl gophermart.Controller) *chi.Mux {
	mux := chi.NewMux()

	mux.Use(middleware.RequestID)
	mux.Use(_middleware.Logger)
	mux.Use(middleware.Recoverer)

	mux.Method(http.MethodPost, "/api/user/register", ctrl.UserRegister())
	mux.Method(http.MethodPost, "/api/user/login", ctrl.UserLogin())

	mux.Route("/api/user/orders", func(r chi.Router) {
		r.Use(_middleware.JWTMiddleware(ctrl.AuthSvc))
		r.Method(http.MethodPost, "/", ctrl.RelateOrderWithUser())
		r.Method(http.MethodGet, "/", ctrl.ListOrdersRelatedWithUser())
	})

	mux.Route("/api/user/balance", func(r chi.Router) {
		r.Use(_middleware.JWTMiddleware(ctrl.AuthSvc))
		r.Method(http.MethodGet, "/", ctrl.BalanceStateByUser())
		r.Method(http.MethodPost, "/withdraw", ctrl.Withdraw())
	})

	mux.Route("/api/user/withdrawals", func(r chi.Router) {
		r.Use(_middleware.JWTMiddleware(ctrl.AuthSvc))
		r.Method(http.MethodGet, "/", ctrl.ListOfWithdrawals())
	})

	return mux
}

func newServer(mux *chi.Mux, addr string) *http.Server {
	return &http.Server{
		Addr:         addr,
		ReadTimeout:  time.Second * 10,
		WriteTimeout: time.Second * 10,
		IdleTimeout:  time.Second * 60,
		Handler:      mux,
	}
}

func run(server *http.Server, sigint chan os.Signal) {
	if err := server.ListenAndServe(); err != nil {
		if errors.Is(err, http.ErrServerClosed) {
			logger.Info("server stopped")
		} else {
			logger.Error("starting failed", "error", err)
		}
	}
	sigint <- os.Interrupt

}

func shutdown(ctx context.Context, server *http.Server) {
	newCtx, cancel := context.WithTimeout(ctx, time.Second*5)
	defer cancel()

	if err := server.Shutdown(newCtx); err != nil {
		logger.Error("server shutdown failed", "error", err)
		server.Close()
	} else {
		logger.Info("server stopped gracefully")
	}
}

func signalSubscription() chan os.Signal {
	s := make(chan os.Signal, 1)
	signal.Notify(s, os.Interrupt)
	return s
}
