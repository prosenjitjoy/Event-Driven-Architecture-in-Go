package main

import (
	"database/sql"
	"log"
	"mall/baskets"
	"mall/customers"
	"mall/depot"
	"mall/internal/config"
	"mall/internal/logger"
	"mall/internal/monolith"
	"mall/internal/waiter"
	"mall/internal/web"
	"mall/notifications"
	"mall/ordering"
	"mall/payments"
	"mall/stores"
	"net/http"

	"github.com/go-chi/chi/v5"
	_ "github.com/lib/pq"
	"github.com/nats-io/nats.go"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	// parse configuration
	cfg, err := config.InitConfig()
	if err != nil {
		return err
	}

	// connect database
	db, err := sql.Open("postgres", cfg.PG.Conn)
	if err != nil {
		return err
	}
	defer db.Close()

	// connect nats
	nc, err := nats.Connect(cfg.Nats.URL)
	if err != nil {
		return err
	}
	defer nc.Close()

	// init jetstream
	js, err := nc.JetStream()
	if err != nil {
		return err
	}

	// init logger
	logger := logger.New(logger.LogConfig{
		Environment: cfg.Environment,
		LogLevel:    cfg.LogLevel,
	})

	// init grpc
	grpcServer := grpc.NewServer()
	reflection.Register(grpcServer)

	// init infrastructure
	m := app{
		cfg:    cfg,
		db:     db,
		nc:     nc,
		js:     js,
		logger: logger,
		mux:    chi.NewMux(),
		rpc:    grpcServer,
		waiter: waiter.New(waiter.CatchSignals()),
		modules: []monolith.Module{
			&baskets.Module{},
			&customers.Module{},
			&depot.Module{},
			&notifications.Module{},
			&ordering.Module{},
			&payments.Module{},
			&stores.Module{},
		},
	}

	// init modules
	if err = m.startupModules(); err != nil {
		return err
	}

	// mount web resources
	m.mux.Mount("/", http.FileServer(http.FS(web.WebUI)))

	log.Println("started mall application")
	defer log.Println("stopped mall application")

	m.waiter.Add(
		m.waitForWeb,
		m.waitForRPC,
		m.waitForStream,
	)

	return m.waiter.Wait()
}
