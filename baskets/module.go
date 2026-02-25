package baskets

import (
	"context"
	"database/sql"
	"log/slog"
	"mall/baskets/basketspb"
	"mall/baskets/internal/application"
	"mall/baskets/internal/constants"
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
	pg "mall/internal/postgres"
	"mall/internal/registry"
	"mall/internal/system"
	"mall/internal/tm"
	"mall/stores/storespb"
)

type Module struct{}

func (*Module) Startup(ctx context.Context, mono system.Service) error {
	return Root(ctx, mono)
}

func Root(ctx context.Context, service system.Service) error {
	container := di.New()

	// setup driven adapters
	container.AddSingleton(constants.RegistryKey, func(c di.Container) (any, error) {
		reg := registry.New()

		if err := domain.Registrations(reg); err != nil {
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

	stream := jetstream.NewStream(service.Config().Nats.Stream, service.JS(), service.Logger())

	container.AddSingleton(constants.DomainDispatcherKey, func(c di.Container) (any, error) {
		return ddd.NewEventDispatcher[ddd.Event](), nil
	})

	container.AddSingleton(constants.DatabaseTransactionKey, func(c di.Container) (any, error) {
		return service.DB().Begin()
	})

	container.AddScoped(constants.MessagePublisherKey, func(c di.Container) (any, error) {
		tx := c.Get(constants.DatabaseTransactionKey).(*sql.Tx)

		outboxStore := pg.NewOutboxStore(constants.OutboxTableName, tx)
		outboxPublisher := tm.OutboxPublisher(outboxStore)

		return am.NewMessagePublisher(
			stream,
			outboxPublisher,
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

	container.AddScoped(constants.BasketsRepoKey, func(c di.Container) (any, error) {
		tx := c.Get(constants.DatabaseTransactionKey).(*sql.Tx)
		reg := c.Get(constants.RegistryKey).(registry.Registry)

		eventStore := pg.NewEventStore(constants.EventsTableName, tx, reg)
		snapshotStore := pg.NewSnapshotStore(constants.SnapshotsTableName, tx, reg)
		aggregateStore := es.AggregateStoreWithMiddleware(eventStore, snapshotStore)

		return es.NewAggregateRepository[*domain.Basket](domain.BasketAggregate, reg, aggregateStore), nil
	})

	container.AddScoped(constants.StoresRepoKey, func(c di.Container) (any, error) {
		tx := c.Get(constants.DatabaseTransactionKey).(*sql.Tx)

		fallback := grpc.NewStoreRepository(service.Config().Rpc.Service(constants.StoresServiceName))

		return postgres.NewStoreCacheRepository(constants.StoresCacheTableName, tx, fallback), nil
	})

	container.AddScoped(constants.ProductsRepoKey, func(c di.Container) (any, error) {
		tx := c.Get(constants.DatabaseTransactionKey).(*sql.Tx)

		fallback := grpc.NewProductRepository(service.Config().Rpc.Service(constants.StoresServiceName))

		return postgres.NewProductCacheRepository(constants.ProductsCacheTableName, tx, fallback), nil
	})

	// setup application
	container.AddScoped(constants.ApplicationKey, func(c di.Container) (any, error) {
		baskets := c.Get(constants.BasketsRepoKey).(domain.BasketRepository)
		stores := c.Get(constants.StoresRepoKey).(domain.StoreCacheRepository)
		products := c.Get(constants.ProductsRepoKey).(domain.ProductCacheRepository)
		domainDispatcher := c.Get(constants.DomainDispatcherKey).(*ddd.EventDispatcher[ddd.Event])

		app := application.New(baskets, stores, products, domainDispatcher)

		return logging.LogApplicationAccess(app, service.Logger()), nil
	})

	container.AddScoped(constants.DomainEventHandlersKey, func(c di.Container) (any, error) {
		eventPublisher := c.Get(constants.EventPublisherKey).(am.EventPublisher)

		domainEventHandlers := handlers.NewDomainEventHandlers(eventPublisher)
		domainEventHandlers = logging.LogEventHandlerAccess[ddd.Event](
			domainEventHandlers,
			"DomainEvents",
			service.Logger(),
		) // logging wrapper

		return domainEventHandlers, nil
	})

	container.AddScoped(constants.IntegrationEventHandlersKey, func(c di.Container) (any, error) {
		reg := di.Get(ctx, constants.RegistryKey).(registry.Registry)
		stores := c.Get(constants.StoresRepoKey).(domain.StoreCacheRepository)
		products := c.Get(constants.ProductsRepoKey).(domain.ProductCacheRepository)

		integrationEventHandlers := handlers.NewIntegrationEventHandlers(stores, products)
		integrationEventHandlers = logging.LogEventHandlerAccess[ddd.Event](
			integrationEventHandlers,
			"IntegrationEvents",
			service.Logger(),
		) // logging wrapper

		inboxStore := di.Get(ctx, constants.InboxStoreKey).(tm.InboxStore)
		inboxHandler := tm.InboxHandler(inboxStore)

		return am.NewEventHandler(reg, integrationEventHandlers, inboxHandler), nil
	})

	outboxStore := pg.NewOutboxStore(constants.OutboxTableName, service.DB())
	outboxProcessor := tm.NewOutboxProcessor(stream, outboxStore)

	// setup driver adapters
	if err := grpc.RegisterServerTx(container, service.RPC()); err != nil {
		return err
	}
	if err := rest.RegisterGateway(ctx, service.Mux(), service.Config().Rpc.Address()); err != nil {
		return err
	}
	if err := rest.RegisterSwagger(service.Mux()); err != nil {
		return err
	}

	handlers.RegisterDomainEventHandlersTx(container)
	err := handlers.RegisterIntegrationEventHandlersTx(container)
	if err != nil {
		return err
	}

	startOutboxProcessor(ctx, outboxProcessor, service.Logger())

	return nil
}

func startOutboxProcessor(ctx context.Context, outboxProcessor tm.OutboxProcessor, logger *slog.Logger) {
	go func() {
		err := outboxProcessor.Start(ctx)
		if err != nil {
			logger.Error("baskets outbox processor encountered an error", "error", err.Error())
		}
	}()
}
