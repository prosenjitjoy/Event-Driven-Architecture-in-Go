package main

import (
	"context"
	"database/sql"
	"log"
	"mall/baskets"
	"mall/cosec"
	"mall/customers"
	"mall/depot"
	"mall/internal/config"
	"mall/internal/system"
	"mall/internal/web"
	"mall/migrations"
	"mall/notifications"
	"mall/ordering"
	"mall/payments"
	"mall/search"
	"mall/stores"
	"net/http"

	"github.com/nats-io/nats.go"
)

type monolith struct {
	*system.System
	modules []system.Module
}

func main() {
	ctx := context.Background()

	if err := run(ctx); err != nil {
		log.Fatalf("monolith service: %s", err)
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

	m := monolith{
		System: s,
		modules: []system.Module{
			&baskets.Module{},
			&customers.Module{},
			&depot.Module{},
			&notifications.Module{},
			&ordering.Module{},
			&payments.Module{},
			&stores.Module{},
			&cosec.Module{},
			&search.Module{},
		},
	}

	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {
			return
		}
	}(m.DB())

	if err := m.MigrateDB(migrations.FS); err != nil {
		return err
	}

	defer func(nc *nats.Conn) {
		nc.Close()
	}(m.NC())

	// init modules
	if err = m.startupModules(); err != nil {
		return err
	}

	// mount web resources
	m.Mux().Mount("/", http.FileServer(http.FS(web.WebUI)))

	m.Logger().InfoContext(ctx, "started mall application")
	defer m.Logger().InfoContext(ctx, "stopped mall application")

	m.Waiter().AddWaitFunc(
		m.WaitForWeb,
		m.WaitForRPC,
		m.WaitForStream,
	)

	return m.Waiter().Wait()
}

func (m *monolith) startupModules() error {
	for _, module := range m.modules {
		ctx := m.Waiter().Context()
		if err := module.Startup(ctx, m); err != nil {
			return err
		}
	}

	return nil
}
