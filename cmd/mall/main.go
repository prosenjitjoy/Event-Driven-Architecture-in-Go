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
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	// parse config
	cfg, err := config.InitConfig()
	if err != nil {
		return err
	}

	m := app{cfg: cfg}

	// init infrastructure
	m.db, err = sql.Open("postgres", cfg.PG.Conn)
	if err != nil {
		return err
	}
	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {
			return
		}
	}(m.db)

	m.logger = logger.New(logger.LogConfig{
		Environment: cfg.Environment,
		LogLevel:    cfg.LogLevel,
	})

	m.rpc = initRpc()
	m.mux = initMux()
	m.waiter = waiter.New(waiter.CatchSignals())

	// init modules
	m.modules = []monolith.Module{
		&baskets.Module{},
		&customers.Module{},
		&depot.Module{},
		&notifications.Module{},
		&ordering.Module{},
		&payments.Module{},
		&stores.Module{},
	}

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
	)

	return m.waiter.Wait()
}

func initRpc() *grpc.Server {
	server := grpc.NewServer()
	reflection.Register(server)

	return server
}

func initMux() *chi.Mux {
	return chi.NewMux()
}
