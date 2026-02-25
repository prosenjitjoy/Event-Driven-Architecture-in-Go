package main

import (
	"context"
	"database/sql"
	"log"
	"mall/baskets"
	"mall/baskets/migrations"
	"mall/internal/config"
	"mall/internal/system"
	"mall/internal/web"
	"net/http"

	_ "github.com/lib/pq"
	"github.com/nats-io/nats.go"
)

func main() {
	ctx := context.Background()

	if err := run(ctx); err != nil {
		log.Fatalf("basket service: %s", err)
	}
}

func run(ctx context.Context) error {
	// parse configuration
	cfg, err := config.InitConfig()
	if err != nil {
		return err
	}

	// add infrastructure
	s, err := system.NewSystem(ctx, cfg)
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
	if err := baskets.Root(s.Waiter().Context(), s); err != nil {
		return err
	}

	s.Logger().InfoContext(ctx, "started baskets service")
	defer s.Logger().InfoContext(ctx, "stopped baskets service")

	s.Waiter().AddWaitFunc(
		s.WaitForWeb,
		s.WaitForRPC,
		s.WaitForStream,
	)

	return s.Waiter().Wait()
}
