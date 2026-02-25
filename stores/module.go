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
	pg "mall/internal/postgres"
	"mall/internal/registry"
	"mall/internal/system"
	"mall/internal/tm"
	"mall/stores/internal/application"
	"mall/stores/internal/constants"
	"mall/stores/internal/domain"
	"mall/stores/internal/grpc"
	"mall/stores/internal/handlers"
	"mall/stores/internal/logging"
	"mall/stores/internal/postgres"
	"mall/stores/internal/rest"
	"mall/stores/storespb"
)

type Module struct{}

func (*Module) Startup(ctx context.Context, mono system.Service) error {
	return Root(ctx, mono)
}

func Root(ctx context.Context, service system.Service) (err error) {
	container := di.New()

	// setup Driven adapters
	container.AddSingleton(constants.RegistryKey, func(c di.Container) (any, error) {
		reg := registry.New()

		if err := domain.Registrations(reg); err != nil {
			return nil, err
		}
		if err := storespb.Registrations(reg); err != nil {
			return nil, err
		}

		return reg, nil
	})

	stream := jetstream.NewStream(service.Config().Nats.Stream, service.JS(), service.Logger())

	container.AddSingleton(constants.DomainDispatcherKey, func(c di.Container) (any, error) {
		return ddd.NewEventDispatcher[ddd.AggregateEvent](), nil
	})

	container.AddScoped(constants.DatabaseTransactionKey, func(c di.Container) (any, error) {
		return service.DB().Begin()
	})

	container.AddScoped(constants.MessagePublisherKey, func(c di.Container) (any, error) {
		tx := c.Get(constants.DatabaseTransactionKey).(*sql.Tx)

		outboxStore := pg.NewOutboxStore(constants.OutboxTableName, tx)
		outboxHandler := tm.OutboxPublisher(outboxStore)

		return am.NewMessagePublisher(
			stream,
			outboxHandler,
		), nil
	})

	container.AddSingleton(constants.MessageSubscriberKey, func(c di.Container) (any, error) {
		return am.NewMessageSubscriber(
			stream,
		), nil
	})

	container.AddScoped(constants.EventPublisherKey, func(c di.Container) (any, error) {
		reg := c.Get(constants.RegistryKey).(registry.Registry)
		msgPublisher := c.Get(constants.MessagePublisherKey).(am.MessagePublisher)

		return am.NewEventPublisher(reg, msgPublisher), nil
	})

	container.AddScoped(constants.InboxStoreKey, func(c di.Container) (any, error) {
		tx := c.Get(constants.DatabaseTransactionKey).(*sql.Tx)

		return pg.NewInboxStore(constants.InboxTableName, tx), nil
	})

	container.AddScoped(constants.AggregateStoreKey, func(c di.Container) (any, error) {
		tx := c.Get(constants.DatabaseTransactionKey).(*sql.Tx)
		reg := c.Get(constants.RegistryKey).(registry.Registry)

		eventStore := pg.NewEventStore(constants.EventsTableName, tx, reg)
		snapshotStore := pg.NewSnapshotStore(constants.SnapshotsTableName, tx, reg)

		return es.AggregateStoreWithMiddleware(eventStore, snapshotStore), nil
	})

	container.AddScoped(constants.StoresRepoKey, func(c di.Container) (any, error) {
		reg := c.Get(constants.RegistryKey).(registry.Registry)
		aggregateStore := c.Get(constants.AggregateStoreKey).(es.AggregateStore)

		return es.NewAggregateRepository[*domain.Store](domain.StoreAggregate, reg, aggregateStore), nil
	})

	container.AddScoped(constants.ProductsRepoKey, func(c di.Container) (any, error) {
		reg := c.Get(constants.RegistryKey).(registry.Registry)
		aggregateStore := c.Get(constants.AggregateStoreKey).(es.AggregateStore)

		return es.NewAggregateRepository[*domain.Product](domain.ProductAggregate, reg, aggregateStore), nil
	})

	container.AddScoped(constants.CatalogRepoKey, func(c di.Container) (any, error) {
		tx := c.Get(constants.DatabaseTransactionKey).(*sql.Tx)

		return postgres.NewCatalogRepository(constants.CatalogTableName, tx), nil
	})

	container.AddScoped(constants.MallRepoKey, func(c di.Container) (any, error) {
		tx := c.Get(constants.DatabaseTransactionKey).(*sql.Tx)

		return postgres.NewMallRepository(constants.MallTableName, tx), nil
	})

	// setup application
	container.AddScoped(constants.ApplicationKey, func(c di.Container) (any, error) {
		stores := c.Get(constants.StoresRepoKey).(domain.StoreRepository)
		products := c.Get(constants.ProductsRepoKey).(domain.ProductRepository)
		catalog := c.Get(constants.CatalogRepoKey).(domain.CatalogRepository)
		mall := c.Get(constants.MallRepoKey).(domain.MallRepository)

		app := application.New(stores, products, catalog, mall)

		return logging.LogApplicationAccess(app, service.Logger()), nil
	})

	container.AddScoped(constants.CatalogHandlersKey, func(c di.Container) (any, error) {
		catalog := c.Get(constants.CatalogRepoKey).(domain.CatalogRepository)

		catalogHandlers := handlers.NewCatalogHandlers(catalog)
		catalogHandlers = logging.LogEventHandlerAccess[ddd.AggregateEvent](
			catalogHandlers,
			"Catalog",
			service.Logger(),
		) // logging wrapper

		return catalogHandlers, nil
	})

	container.AddScoped(constants.MallHandlersKey, func(c di.Container) (any, error) {
		mall := c.Get(constants.MallRepoKey).(domain.MallRepository)

		mallHandlers := handlers.NewMallHandlers(mall)
		mallHandlers = logging.LogEventHandlerAccess[ddd.AggregateEvent](
			mallHandlers,
			"Mall",
			service.Logger(),
		) // logging wrapper

		return mallHandlers, nil
	})

	container.AddScoped(constants.DomainEventHandlersKey, func(c di.Container) (any, error) {
		eventPublisher := c.Get(constants.EventPublisherKey).(am.EventPublisher)

		domainEventHandlers := handlers.NewDomainEventHandlers(eventPublisher)
		domainEventHandlers = logging.LogEventHandlerAccess[ddd.AggregateEvent](
			domainEventHandlers,
			"DomainEvents",
			service.Logger(),
		)

		return domainEventHandlers, nil
	})

	outboxStore := pg.NewOutboxStore(constants.OutboxTableName, service.DB())
	outboxProcessor := tm.NewOutboxProcessor(stream, outboxStore)

	// setup Driver adapters
	if err = grpc.RegisterServerTx(container, service.RPC()); err != nil {
		return err
	}
	if err = rest.RegisterGateway(ctx, service.Mux(), service.Config().Rpc.Address()); err != nil {
		return err
	}
	if err = rest.RegisterSwagger(service.Mux()); err != nil {
		return err
	}

	handlers.RegisterCatalogHandlersTx(container)
	handlers.RegisterMallHandlersTx(container)
	handlers.RegisterDomainEventHandlersTx(container)

	startOutboxProcessor(ctx, outboxProcessor, service.Logger())

	return nil
}

func startOutboxProcessor(ctx context.Context, outboxProcessor tm.OutboxProcessor, logger *slog.Logger) {
	go func() {
		err := outboxProcessor.Start(ctx)
		if err != nil {
			logger.Error("stores outbox processor encountered an error", "error", err.Error())
		}
	}()
}
