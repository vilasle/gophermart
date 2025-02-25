package main

import (
	"context"
	"database/sql"
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

type cliArgs struct {
	addr  string
	dbUrl string
	debug bool
}

func initCli() cliArgs {
	args := cliArgs{}
	pflag.StringVarP(&args.addr, "address", "a", ":8080", "address to listen on")
	pflag.StringVarP(&args.dbUrl, "database", "d", "", "database url e.g postgres://postgres:postgres@localhost:5432/postgres")
	pflag.BoolVarP(&args.debug, "debug", "D", false, "enable debug message")
	pflag.Parse()

	return args
}

func main() {
	args := initCli()

	if args.debug {
		logger.Init(os.Stdout, logger.DebugLevel)
	} else {
		logger.Init(os.Stdout, logger.InfoLevel)
	}

	if args.dbUrl == "" {
		logger.Error("database url is required")
		pflag.Usage()
		os.Exit(1)
	}

	ctx := context.Background()

	em := calculation.NewEventManager(ctx)

	db, err := sql.Open("pgx", args.dbUrl)
	if err != nil {
		logger.Error("Unable to connect to database\n", "url", args.dbUrl)
		os.Exit(1)
	}
	defer db.Close()

	repCalc, err := cRep.NewCalculationRepository(db)
	if err != nil {
		logger.Error("can not init calculation repository", "error", err)
		os.Exit(1)
	}

	ruleSvc := calculation.NewRuleService(calculation.RuleServiceConfig{
		Repository:   repCalc,
		EventManager: em,
	})

	calcSvc := calculation.NewCalculationService(calculation.CalculationServiceConfig{
		CalculationRepository: repCalc,
		CalculationRules:      repCalc,
		EventManager:          em,
	})

	ctrl := accrual.Controller{
		CalculationService:     calcSvc,
		CalculationRuleService: ruleSvc,
	}

	mux := chi.NewMux()

	mux.Use(middleware.RequestID)
	mux.Use(_middleware.Logger)
	mux.Use(middleware.Recoverer)

	mux.Method(http.MethodGet, "/api/orders/{number}", ctrl.OrderInfo())
	mux.Method(http.MethodPost, "/api/orders", ctrl.RegisterOrder())
	mux.Method(http.MethodPost, "/api/goods", ctrl.AddCalculationRules())

	server := http.Server{
		Addr:         ":8080",
		ReadTimeout:  time.Second * 10,
		WriteTimeout: time.Second * 10,
		IdleTimeout:  time.Second * 60,
		Handler:      mux,
	}

	sigint := make(chan os.Signal, 1)
	signal.Notify(sigint, os.Interrupt)

	go func() {
		if err := server.ListenAndServe(); err != nil {
			logger.Error("starting server failed", "error", err)
		}
		sigint <- os.Interrupt
	}()

	<-sigint
	server.Shutdown(ctx)

}
