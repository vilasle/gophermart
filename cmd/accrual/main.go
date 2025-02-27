package main

import (
	"context"
	"database/sql"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/spf13/pflag"

	"github.com/vilasle/gophermart/internal/logger"
	_middleware "github.com/vilasle/gophermart/internal/middleware"

	"github.com/vilasle/gophermart/internal/controller/accrual"
	cRep "github.com/vilasle/gophermart/internal/repository/calculation/postgres"
	"github.com/vilasle/gophermart/internal/service/calculation"
)

func main() {
	args := initCli()

	initLogger(args)

	if err := checkArgs(args); err != nil {
		logger.Error("invalid arguments", "error", err)
		pflag.Usage()
		os.Exit(1)
	}

	db, err := sql.Open("pgx", args.dbURI)
	if err != nil {
		logger.Error("connecting to database failed", "url", args.dbURI, "error", err)
		os.Exit(1)
	}
	defer db.Close()

	repository, err := cRep.NewCalculationRepository(db)
	if err != nil {
		logger.Error("can not init calculation repository", "error", err)
		os.Exit(1)
	}

	ctx := context.Background()

	em := calculation.NewEventManager(ctx)
	defer em.Stop()

	mux := newMux(newController(repository, em))

	server := newServer(mux, args.addr)
	defer server.Close()

	s := signalSubscription()

	go run(server, s)

	<-s
}

type cliArgs struct {
	addr  string
	dbURI string
	debug bool
}

func initCli() cliArgs {
	args := cliArgs{}
	pflag.StringVarP(&args.addr, "address", "a", ":8080", "address to listen on")
	pflag.StringVarP(&args.dbURI, "database", "d", "", "database url e.g postgres://postgres:postgres@localhost:5432/postgres")
	pflag.BoolVarP(&args.debug, "debug", "D", false, "enable debug message")
	pflag.Parse()

	args.addr = getEnv("RUN_ADDRESS", args.addr)
	args.dbURI = getEnv("DATABASE_URI", args.dbURI)

	return args
}

func getEnv(key, fallback string) string {
	result := fallback
	if value, _ := os.LookupEnv(key); value != "" {
		result = value
	}
	return result
}

func initLogger(args cliArgs) {
	if args.debug {
		logger.Init(os.Stdout, logger.DebugLevel)
	} else {
		logger.Init(os.Stdout, logger.InfoLevel)
	}
}

func checkArgs(args cliArgs) error {
	errs := make([]error, 2)
	if args.dbURI == "" {
		errs = append(errs, errors.New("database url is required"))
	}

	if args.addr == "" {
		errs = append(errs, errors.New("address is required"))
	}
	return errors.Join(errs...)
}

func newController(repository cRep.CalculationRepository, eventManager *calculation.EventManager) accrual.Controller {
	ruleSvc := calculation.NewRuleService(calculation.RuleServiceConfig{
		Repository:   repository,
		EventManager: eventManager,
	})

	calcSvc := calculation.NewCalculationService(calculation.CalculationServiceConfig{
		CalculationRepository: repository,
		CalculationRules:      repository,
		EventManager:          eventManager,
	})

	return accrual.Controller{
		CalculationService:     calcSvc,
		CalculationRuleService: ruleSvc,
	}
}

func newMux(ctrl accrual.Controller) *chi.Mux {
	mux := chi.NewMux()

	mux.Use(middleware.RequestID)
	mux.Use(_middleware.Logger)
	mux.Use(middleware.Recoverer)

	mux.Method(http.MethodGet, "/api/orders/{number}", ctrl.OrderInfo())
	mux.Method(http.MethodGet, "/orders/{number}", ctrl.OrderInfo())
	mux.Method(http.MethodPost, "/api/orders", ctrl.RegisterOrder())
	mux.Method(http.MethodPost, "/api/goods", ctrl.AddCalculationRules())

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

func signalSubscription() chan os.Signal {
	s := make(chan os.Signal, 1)
	signal.Notify(s, os.Interrupt)
	return s
}
