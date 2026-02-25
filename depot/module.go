package depot

import (
	"context"
	"database/sql"
	"log/slog"

	"mall/depot/depotpb"
	"mall/depot/internal/application"
	"mall/depot/internal/constants"
	"mall/depot/internal/domain"
	"mall/depot/internal/grpc"
	"mall/depot/internal/handlers"
	"mall/depot/internal/logging"
	"mall/depot/internal/postgres"
	"mall/depot/internal/rest"
	"mall/internal/am"
	"mall/internal/ddd"
	"mall/internal/di"
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

func Root(ctx context.Context, service system.Service) (err error) {
	container := di.New()

	// setup Driven adapters
	container.AddSingleton(constants.RegistryKey, func(c di.Container) (any, error) {
		reg := registry.New()
		if err := storespb.Registrations(reg); err != nil {
			return nil, err
		}
		if err := depotpb.Registrations(reg); err != nil {
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

	container.AddScoped(constants.CommandPublisherKey, func(c di.Container) (any, error) {
		reg := c.Get(constants.RegistryKey).(registry.Registry)
		msgPublisher := c.Get(constants.MessagePublisherKey).(am.MessagePublisher)

		return am.NewCommandPublisher(reg, msgPublisher), nil
	})

	container.AddScoped(constants.ReplyPublisherKey, func(c di.Container) (any, error) {
		reg := c.Get(constants.RegistryKey).(registry.Registry)
		msgPublisher := c.Get(constants.MessagePublisherKey).(am.MessagePublisher)

		return am.NewReplyPublisher(reg, msgPublisher), nil
	})

	container.AddScoped(constants.InboxStoreKey, func(c di.Container) (any, error) {
		tx := c.Get(constants.DatabaseTransactionKey).(*sql.Tx)

		return pg.NewInboxStore(constants.InboxTableName, tx), nil
	})

	container.AddScoped(constants.ShoppingListsRepoKey, func(c di.Container) (any, error) {
		tx := c.Get(constants.DatabaseTransactionKey).(*sql.Tx)

		return postgres.NewShoppingListRepository(constants.ShoppingListsTableName, tx), nil
	})

	container.AddScoped(constants.StoresCacheRepoKey, func(c di.Container) (any, error) {
		tx := c.Get(constants.DatabaseTransactionKey).(*sql.Tx)

		fallback := grpc.NewStoreRepository(service.Config().Rpc.Address())

		return postgres.NewStoreCacheRepository(constants.StoresCacheTableName, tx, fallback), nil
	})

	container.AddScoped(constants.ProductsCacheRepoKey, func(c di.Container) (any, error) {
		tx := c.Get(constants.DatabaseTransactionKey).(*sql.Tx)

		fallback := grpc.NewProductRepository(service.Config().Rpc.Address())

		return postgres.NewProductCacheRepository(constants.ProductsCacheTableName, tx, fallback), nil
	})

	// setup application
	container.AddScoped(constants.ApplicationKey, func(c di.Container) (any, error) {
		shoppingList := c.Get(constants.ShoppingListsRepoKey).(domain.ShoppingListRepository)
		storeCache := c.Get(constants.StoresCacheRepoKey).(domain.StoreCacheRepository)
		productCache := c.Get(constants.ProductsCacheRepoKey).(domain.ProductCacheRepository)
		domainDispatcher := c.Get(constants.DomainDispatcherKey).(*ddd.EventDispatcher[ddd.AggregateEvent])

		app := application.New(shoppingList, storeCache, productCache, domainDispatcher)

		return logging.LogApplicationAccess(app, service.Logger()), nil
	})

	container.AddScoped(constants.DomainEventHandlersKey, func(c di.Container) (any, error) {
		eventPublisher := c.Get(constants.EventPublisherKey).(am.EventPublisher)

		domainEventHandlers := handlers.NewDomainEventHandlers(eventPublisher)
		domainEventHandlers = logging.LogEventHandlerAccess[ddd.AggregateEvent](
			domainEventHandlers,
			"DomainEvents",
			service.Logger(),
		) // logging wrapper

		return domainEventHandlers, nil
	})

	container.AddScoped(constants.IntegrationEventHandlersKey, func(c di.Container) (any, error) {
		return logging.LogEventHandlerAccess[ddd.Event](
			handlers.NewIntegrationEventHandlers(
				c.Get("stores").(domain.StoreCacheRepository),
				c.Get("products").(domain.ProductCacheRepository),
			),
			"IntegrationEvents", c.Get("logger").(*slog.Logger),
		), nil
	})

	container.AddScoped(constants.CommandHandlersKey, func(c di.Container) (any, error) {
		reg := c.Get(constants.RegistryKey).(registry.Registry)
		app := c.Get(constants.ApplicationKey).(application.App)

		commandHandlers := handlers.NewCommandHandlers(app)
		commandHandlers = logging.LogCommandHandlerAccess(
			commandHandlers,
			"Commands",
			service.Logger(),
		) // logging wrapper

		replyPublisher := c.Get(constants.ReplyPublisherKey).(am.ReplyPublisher)
		inboxStore := c.Get(constants.InboxStoreKey).(tm.InboxStore)
		inboxHandler := tm.InboxHandler(inboxStore)

		return am.NewCommandHandler(reg, replyPublisher, commandHandlers, inboxHandler), nil
	})

	outboxStore := pg.NewOutboxStore(constants.OutboxTableName, service.DB())
	outboxProcessor := tm.NewOutboxProcessor(stream, outboxStore)

	// setup Driver adapters
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

	if err = handlers.RegisterIntegrationEventHandlersTx(container); err != nil {
		return err
	}
	if err = handlers.RegisterCommandHandlersTx(container); err != nil {
		return err
	}

	startOutboxProcessor(ctx, outboxProcessor, service.Logger())

	return nil
}

func startOutboxProcessor(ctx context.Context, outboxProcessor tm.OutboxProcessor, logger *slog.Logger) {
	go func() {
		err := outboxProcessor.Start(ctx)
		if err != nil {
			logger.Error("depot outbox processor encountered an error", "error", err.Error())
		}
	}()
}
