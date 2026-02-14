package main

import (
	"database/sql"
	"log"

	"mall/internal/config"
	"mall/internal/system"
	"mall/internal/web"
	"mall/payments"
	"mall/payments/migrations"
	"net/http"

	_ "github.com/lib/pq"
	"github.com/nats-io/nats.go"
)

func main() {
	if err := run(); err != nil {
		log.Fatalf("payments service: %s", err)
	}
}

func run() error {
	// parse configuration
	cfg, err := config.InitConfig()
	if err != nil {
		return err
	}

	// add infrastructure
	s, err := system.NewSystem(cfg)
	if err != nil {
		return err
	}

	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {
			return
		}
	}(s.DB())

	if err := s.MigrateDB(migrations.FS); err != nil {
		return err
	}

	defer func(nc *nats.Conn) {
		nc.Close()
	}(s.NC())

	// mount web resources
	s.Mux().Mount("/", http.FileServer(http.FS(web.WebUI)))

	// call the module composition root
	if err := payments.Root(s.Waiter().Context(), s); err != nil {
		return err
	}

	s.Logger().Info("started payments service")
	defer s.Logger().Info("stopped payments service")

	s.Waiter().Add(
		s.WaitForWeb,
		s.WaitForRPC,
		s.WaitForStream,
	)

	return s.Waiter().Wait()
}
