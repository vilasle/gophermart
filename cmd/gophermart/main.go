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

	// if args.debug {
	logger.Init(os.Stdout, logger.DebugLevel)
	// } else {
	// 	logger.Init(os.Stdout, logger.InfoLevel)
	// }

	if args.dbURI == "" {
		logger.Error("database url is required")
		pflag.Usage()
		os.Exit(1)
	}

	accrualUrl, err := url.Parse(args.accrualAddr)
	if err != nil {
		logger.Error("can not parse accrual endpoint address", "error", err)
		os.Exit(1)
	}

	ctx := context.Background()

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

	withdrawalSvc := withdrawal.NewWithdrawalService(dbRep)

	authSvc := authorization.NewAuthorizationService(dbRep)

	accrualSvc := accrual.NewAccrualService(
		httpRep.NewAccrualRepository(accrualUrl),
	)

	orderSvc := order.NewOrderService(dbRep, accrualSvc, dbRep)

	ctrl := gophermart.Controller{
		AuthSvc:     authSvc,
		OrderSvc:    orderSvc,
		WithdrawSvc: withdrawalSvc,
	}

	mux := chi.NewMux()

	mux.Use(middleware.RequestID)
	mux.Use(_middleware.Logger)
	mux.Use(middleware.Recoverer)

	mux.Method(http.MethodPost, "/api/user/register", ctrl.UserRegister())
	mux.Method(http.MethodPost, "/api/user/login", ctrl.UserLogin())

	mux.Route("/api/user/orders", func(r chi.Router) {
		r.Use(_middleware.JWTMiddleware(authSvc))
		r.Method(http.MethodPost, "/", ctrl.RelateOrderWithUser())
		r.Method(http.MethodGet, "/", ctrl.ListOrdersRelatedWithUser())
	})

	mux.Route("/api/user/balance", func(r chi.Router) {
		r.Use(_middleware.JWTMiddleware(authSvc))
		r.Method(http.MethodGet, "/", ctrl.BalanceStateByUser())
		r.Method(http.MethodPost, "/withdraw", ctrl.Withdraw())
	})

	mux.Route("/api/user/withdrawals", func(r chi.Router) {
		r.Use(_middleware.JWTMiddleware(authSvc))
		r.Method(http.MethodGet, "/", ctrl.ListOfWithdrawals())
	})

	server := http.Server{
		Addr:         args.addr,
		ReadTimeout:  time.Second * 10,
		WriteTimeout: time.Second * 10,
		IdleTimeout:  time.Second * 60,
		Handler:      mux,
	}

	sigint := make(chan os.Signal, 1)
	signal.Notify(sigint, os.Interrupt)

	go func() {
		logger.Info("starting server", "address", args.addr)
		if err := server.ListenAndServe(); err != nil {
			if errors.Is(err, http.ErrServerClosed) {
				logger.Info("starting stopped")
			}
			logger.Error("starting failed", "error", err)
		}
		sigint <- os.Interrupt
	}()

	<-sigint
	server.Shutdown(ctx)

}
