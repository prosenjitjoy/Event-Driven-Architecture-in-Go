package search

import (
	"context"
	"database/sql"
	"log/slog"
	"mall/customers/customerspb"
	"mall/internal/di"
	"mall/internal/jetstream"
	"mall/internal/monolith"
	pg "mall/internal/postgres"
	"mall/internal/registry"
	"mall/internal/tm"
	"mall/ordering/orderingpb"
	"mall/search/internal/application"
	"mall/search/internal/domain"
	"mall/search/internal/grpc"
	"mall/search/internal/handlers"
	"mall/search/internal/logging"
	"mall/search/internal/postgres"
	"mall/search/internal/rest"
	"mall/stores/storespb"
)

type Module struct{}

func (*Module) Startup(ctx context.Context, mono monolith.Monolith) error {
	container := di.New()

	// setup driven adapters
	container.AddSingleton("registry", func(c di.Container) (any, error) {
		reg := registry.New()

		if err := orderingpb.Registrations(reg); err != nil {
			return nil, err
		}
		if err := customerspb.Registrations(reg); err != nil {
			return nil, err
		}
		if err := storespb.Registrations(reg); err != nil {
			return nil, err
		}

		return reg, nil
	})

	container.AddSingleton("logger", func(c di.Container) (any, error) {
		return mono.Logger(), nil
	})

	container.AddSingleton("stream", func(c di.Container) (any, error) {
		logger := c.Get("logger").(*slog.Logger)

		return jetstream.NewStream(mono.Config().Nats.Stream, mono.JS(), logger), nil
	})

	container.AddSingleton("db", func(c di.Container) (any, error) {
		return mono.DB(), nil
	})

	container.AddSingleton("conn", func(c di.Container) (any, error) {
		return grpc.Dial(ctx, mono.Config().Rpc.Address())
	})

	container.AddScoped("tx", func(c di.Container) (any, error) {
		db := c.Get("db").(*sql.DB)

		return db.Begin()
	})

	container.AddScoped("inboxMiddleware", func(c di.Container) (any, error) {
		tx := c.Get("tx").(*sql.Tx)

		inboxStore := pg.NewInboxStore("search.inbox", tx)

		return tm.NewInboxHandlerMiddleware(inboxStore), nil
	})

	container.AddScoped("customers", func(c di.Container) (any, error) {
		tx := c.Get("tx").(*sql.Tx)
		conn := c.Get("conn").(*grpc.ClientConn)

		fallback := grpc.NewCustomerRepository(conn)

		return postgres.NewCustomerCacheRepository("search.customers.cache", tx, fallback), nil
	})

	container.AddScoped("stores", func(c di.Container) (any, error) {
		tx := c.Get("tx").(*sql.Tx)
		conn := c.Get("conn").(*grpc.ClientConn)

		fallback := grpc.NewStoreRepository(conn)

		return postgres.NewStoreCacheRepository("search.stores_cache", tx, fallback), nil
	})

	container.AddScoped("products", func(c di.Container) (any, error) {
		tx := c.Get("tx").(*sql.Tx)
		conn := c.Get("conn").(*grpc.ClientConn)

		fallback := grpc.NewProductRepository(conn)

		return postgres.NewProductCacheRepository("search.products_cache", tx, fallback), nil
	})

	container.AddScoped("orders", func(c di.Container) (any, error) {
		tx := c.Get("tx").(*sql.Tx)

		return postgres.NewOrderRepository("search.orders", tx), nil
	})

	// setup application
	container.AddScoped("app", func(c di.Container) (any, error) {
		logger := c.Get("logger").(*slog.Logger)
		orders := c.Get("orders").(domain.OrderRepository)

		app := application.New(orders)

		return logging.LogApplicationAccess(app, logger), nil
	})

	container.AddScoped("integrationEventHandlers", func(c di.Container) (any, error) {
		logger := c.Get("logger").(*slog.Logger)
		orders := c.Get("orders").(domain.OrderRepository)
		customers := c.Get("customers").(domain.CustomerCacheRepository)
		stores := c.Get("stores").(domain.StoreCacheRepository)
		products := c.Get("products").(domain.ProductCacheRepository)

		integrationEventHandlers := handlers.NewIntegrationHandlers(orders, customers, stores, products)

		return logging.LogEventHandlerAccess(integrationEventHandlers, "IntegrationEvents", logger), nil
	})

	// setup driver adapters
	if err := grpc.RegisterServerTx(container, mono.RPC()); err != nil {
		return err
	}
	if err := rest.RegisterGateway(ctx, mono.Mux(), mono.Config().Rpc.Address()); err != nil {
		return err
	}
	if err := rest.RegisterSwagger(mono.Mux()); err != nil {
		return err
	}

	if err := handlers.RegisterIntegrationEventHandlersTx(container); err != nil {
		return err
	}

	return nil
}
