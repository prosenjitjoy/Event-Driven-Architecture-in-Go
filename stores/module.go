package stores

import (
	"context"
	"database/sql"
	"log/slog"

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
	"mall/stores/internal/application"
	"mall/stores/internal/domain"
	"mall/stores/internal/grpc"
	"mall/stores/internal/handlers"
	"mall/stores/internal/logging"
	"mall/stores/internal/postgres"
	"mall/stores/internal/rest"
	"mall/stores/storespb"
)

type Module struct {
}

func (m *Module) Startup(ctx context.Context, mono monolith.Monolith) (err error) {
	container := di.New()

	// setup Driven adapters
	container.AddSingleton("registry", func(c di.Container) (any, error) {
		reg := registry.New()

		if err := registrations(reg); err != nil {
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
		return ddd.NewEventDispatcher[ddd.AggregateEvent](), nil
	})

	container.AddSingleton("db", func(c di.Container) (any, error) {
		return mono.DB(), nil
	})

	container.AddSingleton("outboxProcessor", func(c di.Container) (any, error) {
		db := c.Get("db").(*sql.DB)
		stream := c.Get("stream").(am.RawMessageStream)

		outboxStore := pg.NewOutboxStore("stores.outbox", db)

		return tm.NewOutboxProcessor(stream, outboxStore), nil
	})

	container.AddScoped("tx", func(c di.Container) (any, error) {
		db := c.Get("db").(*sql.DB)

		return db.Begin()
	})

	container.AddScoped("txStream", func(c di.Container) (any, error) {
		tx := c.Get("tx").(*sql.Tx)
		stream := c.Get("stream").(am.RawMessageStream)

		outboxStore := pg.NewOutboxStore("stores.outbox", tx)
		outboxStream := tm.NewOutboxStreamMiddleware(outboxStore)

		return am.RawMessageStreamWithMiddleware(stream, outboxStream), nil
	})

	container.AddScoped("eventStream", func(c di.Container) (any, error) {
		reg := c.Get("registry").(registry.Registry)
		stream := c.Get("txStream").(am.RawMessageStream)

		return am.NewEventStream(reg, stream), nil
	})

	container.AddScoped("aggregateStore", func(c di.Container) (any, error) {
		tx := c.Get("tx").(*sql.Tx)
		reg := c.Get("registry").(registry.Registry)
		domainDispatcher := c.Get("domainDispatcher").(*ddd.EventDispatcher[ddd.AggregateEvent])

		eventStore := pg.NewEventStore("stores.events", tx, reg)
		eventPublisher := es.NewEventPublisher(domainDispatcher)
		snapshotStore := pg.NewSnapshotStore("stores.snapshots", tx, reg)

		return es.AggregateStoreWithMiddleware(eventStore, eventPublisher, snapshotStore), nil
	})

	container.AddScoped("stores", func(c di.Container) (any, error) {
		reg := c.Get("registry").(registry.Registry)
		aggregateStore := c.Get("aggregateStore").(es.AggregateStore)

		return es.NewAggregateRepository[*domain.Store](domain.StoreAggregate, reg, aggregateStore), nil
	})

	container.AddScoped("products", func(c di.Container) (any, error) {
		reg := c.Get("registry").(registry.Registry)
		aggregateStore := c.Get("aggregateStore").(es.AggregateStore)

		return es.NewAggregateRepository[*domain.Product](domain.ProductAggregate, reg, aggregateStore), nil
	})

	container.AddScoped("catalog", func(c di.Container) (any, error) {
		tx := c.Get("tx").(*sql.Tx)

		return postgres.NewCatalogRepository("stores.products", tx), nil
	})
	container.AddScoped("mall", func(c di.Container) (any, error) {
		tx := c.Get("tx").(*sql.Tx)

		return postgres.NewMallRepository("stores.stores", tx), nil
	})

	// setup application
	container.AddScoped("app", func(c di.Container) (any, error) {
		logger := c.Get("logger").(*slog.Logger)
		stores := c.Get("stores").(domain.StoreRepository)
		products := c.Get("products").(domain.ProductRepository)
		catalog := c.Get("catalog").(domain.CatalogRepository)
		mall := c.Get("mall").(domain.MallRepository)

		app := application.New(stores, products, catalog, mall)

		return logging.LogApplicationAccess(app, logger), nil
	})

	container.AddScoped("catalogHandlers", func(c di.Container) (any, error) {
		logger := c.Get("logger").(*slog.Logger)
		catalog := c.Get("catalog").(domain.CatalogRepository)

		catalogHandlers := handlers.NewCatalogHandlers(catalog)

		return logging.LogEventHandlerAccess[ddd.AggregateEvent](catalogHandlers, "Catalog", logger), nil
	})

	container.AddScoped("mallHandlers", func(c di.Container) (any, error) {
		logger := c.Get("logger").(*slog.Logger)
		mall := c.Get("mall").(domain.MallRepository)

		mallHandlers := handlers.NewMallHandlers(mall)

		return logging.LogEventHandlerAccess[ddd.AggregateEvent](mallHandlers, "Mall", logger), nil
	})

	container.AddScoped("domainEventHandlers", func(c di.Container) (any, error) {
		logger := c.Get("logger").(*slog.Logger)
		eventStream := c.Get("eventStream").(am.EventStream)

		domainEventHandlers := handlers.NewDomainEventHandlers(eventStream)

		return logging.LogEventHandlerAccess[ddd.AggregateEvent](domainEventHandlers, "DomainEvents", logger), nil
	})

	// setup Driver adapters
	if err = grpc.RegisterServerTx(container, mono.RPC()); err != nil {
		return err
	}
	if err = rest.RegisterGateway(ctx, mono.Mux(), mono.Config().Rpc.Address()); err != nil {
		return err
	}
	if err = rest.RegisterSwagger(mono.Mux()); err != nil {
		return err
	}

	handlers.RegisterCatalogHandlersTx(container)
	handlers.RegisterMallHandlersTx(container)
	handlers.RegisterDomainEventHandlersTx(container)

	startOutboxProcessor(ctx, container)

	return nil
}

func registrations(reg registry.Registry) (err error) {
	serde := serdes.NewJsonSerde(reg)

	// Store
	if err = serde.Register(domain.Store{}, func(v any) error {
		store := v.(*domain.Store)
		store.Aggregate = es.NewAggregate("", domain.StoreAggregate)
		return nil
	}); err != nil {
		return
	}
	// store events
	if err = serde.Register(domain.StoreCreated{}); err != nil {
		return
	}
	if err = serde.RegisterKey(domain.StoreParticipationEnabledEvent, domain.StoreParticipationToggled{}); err != nil {
		return
	}
	if err = serde.RegisterKey(domain.StoreParticipationDisabledEvent, domain.StoreParticipationToggled{}); err != nil {
		return
	}
	if err = serde.Register(domain.StoreRebranded{}); err != nil {
		return
	}
	// store snapshots
	if err = serde.RegisterKey(domain.StoreV1{}.SnapshotName(), domain.StoreV1{}); err != nil {
		return
	}

	// Product
	if err = serde.Register(domain.Product{}, func(v any) error {
		store := v.(*domain.Product)
		store.Aggregate = es.NewAggregate("", domain.ProductAggregate)
		return nil
	}); err != nil {
		return
	}
	// product events
	if err = serde.Register(domain.ProductAdded{}); err != nil {
		return
	}
	if err = serde.Register(domain.ProductRebranded{}); err != nil {
		return
	}
	if err = serde.RegisterKey(domain.ProductPriceIncreasedEvent, domain.ProductPriceChanged{}); err != nil {
		return
	}
	if err = serde.RegisterKey(domain.ProductPriceDecreasedEvent, domain.ProductPriceChanged{}); err != nil {
		return
	}
	if err = serde.Register(domain.ProductRemoved{}); err != nil {
		return
	}
	// product snapshots
	if err = serde.RegisterKey(domain.ProductV1{}.SnapshotName(), domain.ProductV1{}); err != nil {
		return
	}

	return
}
func startOutboxProcessor(ctx context.Context, container di.Container) {
	outboxProcessor := container.Get("outboxProcessor").(tm.OutboxProcessor)
	logger := container.Get("logger").(*slog.Logger)

	go func() {
		err := outboxProcessor.Start(ctx)
		if err != nil {
			logger.Error("stores outbox processor encountered an error", "error", err.Error())
		}
	}()
}
