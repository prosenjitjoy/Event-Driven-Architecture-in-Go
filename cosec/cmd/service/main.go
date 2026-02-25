package main

import (
	"context"
	"database/sql"
	"log"
	"mall/cosec"
	"mall/cosec/migrations"
	"mall/internal/config"
	"mall/internal/system"
	"mall/internal/web"
	"net/http"

	"github.com/nats-io/nats.go"
)

func main() {
	if err := run(); err != nil {
		log.Fatalf("cosec service: %s", err)
	}
}

func run() error {
	// parse configuration
	cfg, err := config.InitConfig()
	if err != nil {
		return err
	}

	// add infrastructure
	s, err := system.NewSystem(context.Background(), cfg)
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
	if err := cosec.Root(s.Waiter().Context(), s); err != nil {
		return err
	}

	s.Logger().Info("started cosec service")
	defer s.Logger().Info("stopped cosec service")

	s.Waiter().AddWaitFunc(
		s.WaitForWeb,
		s.WaitForRPC,
		s.WaitForStream,
	)

	return s.Waiter().Wait()
}
