package baskets

import (
	"context"
	"database/sql"
	"log/slog"
	"mall/baskets/basketspb"
	"mall/baskets/internal/application"
	"mall/baskets/internal/domain"
	"mall/baskets/internal/grpc"
	"mall/baskets/internal/handlers"
	"mall/baskets/internal/logging"
	"mall/baskets/internal/postgres"
	"mall/baskets/internal/rest"
	"mall/internal/am"
	"mall/internal/ddd"
	"mall/internal/di"
	"mall/internal/es"
	"mall/internal/jetstream"
	"mall/internal/monolith"
	pg "mall/internal/postgres"
	"mall/internal/registry"
	"mall/internal/registry/serdes"
	"mall/internal/tm"
	"mall/stores/storespb"
)

type Module struct{}

func (*Module) Startup(ctx context.Context, mono monolith.Monolith) error {
	container := di.New()

	// setup driven adapters
	container.AddSingleton("registry", func(c di.Container) (any, error) {
		reg := registry.New()

		if err := registrations(reg); err != nil {
			return nil, err
		}
		if err := basketspb.Registrations(reg); err != nil {
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

	container.AddSingleton("domainDispatcher", func(c di.Container) (any, error) {
		return ddd.NewEventDispatcher[ddd.Event](), nil
	})

	container.AddSingleton("db", func(c di.Container) (any, error) {
		return mono.DB(), nil
	})

	container.AddSingleton("conn", func(c di.Container) (any, error) {
		return grpc.Dial(ctx, mono.Config().Rpc.Address())
	})

	container.AddSingleton("outboxProcessor", func(c di.Container) (any, error) {
		db := c.Get("db").(*sql.DB)
		stream := c.Get("stream").(am.RawMessageStream)
		outboxStore := pg.NewOutboxStore("baskets.outbox", db)

		return tm.NewOutboxProcessor(stream, outboxStore), nil
	})

	container.AddScoped("tx", func(c di.Container) (any, error) {
		db := c.Get("db").(*sql.DB)

		return db.Begin()
	})

	container.AddScoped("txStream", func(c di.Container) (any, error) {
		tx := c.Get("tx").(*sql.Tx)
		outboxStore := pg.NewOutboxStore("baskets.outbox", tx)

		stream := c.Get("stream").(am.RawMessageStream)
		outboxStream := tm.NewOutboxStreamMiddleware(outboxStore)

		return am.RawMessageStreamWithMiddleware(stream, outboxStream), nil
	})

	container.AddScoped("eventStream", func(c di.Container) (any, error) {
		reg := c.Get("registry").(registry.Registry)
		stream := c.Get("txStream").(am.RawMessageStream)

		return am.NewEventStream(reg, stream), nil
	})

	container.AddScoped("inboxMiddleware", func(c di.Container) (any, error) {
		tx := c.Get("tx").(*sql.Tx)
		inboxStore := pg.NewInboxStore("baskets.inbox", tx)

		return tm.NewInboxHandlerMiddleware(inboxStore), nil
	})

	container.AddScoped("aggregateStore", func(c di.Container) (any, error) {
		tx := c.Get("tx").(*sql.Tx)
		reg := c.Get("registry").(registry.Registry)

		eventStore := pg.NewEventStore("baskets.events", tx, reg)
		snapshotStore := pg.NewSnapshotStore("baskets.snapshots", tx, reg)

		return es.AggregateStoreWithMiddleware(eventStore, snapshotStore), nil
	})

	container.AddScoped("baskets", func(c di.Container) (any, error) {
		reg := c.Get("registry").(registry.Registry)
		aggregateStore := c.Get("aggregateStore").(es.AggregateStore)

		return es.NewAggregateRepository[*domain.Basket](domain.BasketAggregate, reg, aggregateStore), nil
	})

	container.AddScoped("stores", func(c di.Container) (any, error) {
		tx := c.Get("tx").(*sql.Tx)
		conn := c.Get("conn").(*grpc.ClientConn)
		fallback := grpc.NewStoreRepository(conn)

		return postgres.NewStoreCacheRepository("baskets.stores_cache", tx, fallback), nil
	})

	container.AddScoped("products", func(c di.Container) (any, error) {
		tx := c.Get("tx").(*sql.Tx)
		conn := c.Get("conn").(*grpc.ClientConn)
		fallback := grpc.NewProductRepository(conn)

		return postgres.NewProductCacheRepository("baskets.products_cache", tx, fallback), nil
	})

	// setup application
	container.AddScoped("app", func(c di.Container) (any, error) {
		logger := c.Get("logger").(*slog.Logger)
		baskets := c.Get("baskets").(domain.BasketRepository)
		stores := c.Get("stores").(domain.StoreCacheRepository)
		products := c.Get("products").(domain.ProductCacheRepository)
		domainDispatcher := c.Get("domainDispatcher").(*ddd.EventDispatcher[ddd.Event])

		app := application.New(baskets, stores, products, domainDispatcher)

		return logging.LogApplicationAccess(app, logger), nil
	})

	container.AddScoped("domainEventHandlers", func(c di.Container) (any, error) {
		logger := c.Get("logger").(*slog.Logger)
		eventStream := c.Get("eventStream").(am.EventStream)

		domainEventHandlers := handlers.NewDomainEventHandlers(eventStream)

		return logging.LogEventHandlerAccess[ddd.Event](domainEventHandlers, "DomainEvents", logger), nil
	})

	container.AddScoped("integrationEventHandlers", func(c di.Container) (any, error) {
		logger := c.Get("logger").(*slog.Logger)

		stores := c.Get("stores").(domain.StoreCacheRepository)
		products := c.Get("products").(domain.ProductCacheRepository)

		integrationEventHandlers := handlers.NewIntegrationEventHandlers(stores, products)

		return logging.LogEventHandlerAccess[ddd.Event](integrationEventHandlers, "IntegrationEvents", logger), nil
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

	handlers.RegisterDomainEventHandlersTx(container)
	err := handlers.RegisterIntegrationEventHandlersTx(container)
	if err != nil {
		return err
	}

	startOutboxProcessor(ctx, container)

	return nil
}

func registrations(reg registry.Registry) error {
	serde := serdes.NewJsonSerde(reg)

	// Basket
	if err := serde.Register(domain.Basket{}, func(v any) error {
		basket := v.(*domain.Basket)
		basket.Items = make(map[string]domain.Item)
		return nil
	}); err != nil {
		return err
	}

	// basket events
	if err := serde.Register(domain.BasketStarted{}); err != nil {
		return err
	}
	if err := serde.Register(domain.BasketCanceled{}); err != nil {
		return err
	}
	if err := serde.Register(domain.BasketCheckedOut{}); err != nil {
		return err
	}
	if err := serde.Register(domain.BasketItemAdded{}); err != nil {
		return err
	}
	if err := serde.Register(domain.BasketItemRemoved{}); err != nil {
		return err
	}

	// basket snapshots
	if err := serde.RegisterKey(domain.BasketV1{}.SnapshotName(), domain.BasketV1{}); err != nil {
		return err
	}

	return nil
}

func startOutboxProcessor(ctx context.Context, container di.Container) {
	logger := container.Get("logger").(*slog.Logger)
	outboxProcessor := container.Get("outboxProcessor").(tm.OutboxProcessor)

	go func() {
		err := outboxProcessor.Start(ctx)
		if err != nil {
			logger.Error("baskets outbox processor encountered an error", "error", err.Error())
		}
	}()
}
