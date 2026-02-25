package ordering

import (
	"context"
	"database/sql"
	"log/slog"
	"mall/baskets/basketspb"
	"mall/depot/depotpb"
	"mall/internal/am"
	"mall/internal/ddd"
	"mall/internal/di"
	"mall/internal/es"
	"mall/internal/jetstream"
	pg "mall/internal/postgres"
	"mall/internal/registry"
	"mall/internal/system"
	"mall/internal/tm"
	"mall/ordering/internal/application"
	"mall/ordering/internal/constants"
	"mall/ordering/internal/domain"
	"mall/ordering/internal/grpc"
	"mall/ordering/internal/handlers"
	"mall/ordering/internal/logging"
	"mall/ordering/internal/rest"
	"mall/ordering/orderingpb"
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
		if err := orderingpb.Registrations(reg); err != nil {
			return nil, err
		}
		if err := depotpb.Registrations(reg); err != nil {
			return nil, err
		}

		return reg, nil
	})

	stream := jetstream.NewStream(service.Config().Nats.Stream, service.JS(), service.Logger())

	container.AddSingleton(constants.DomainDispatcherKey, func(c di.Container) (any, error) {
		return ddd.NewEventDispatcher[ddd.Event](), nil
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

	container.AddScoped(constants.MessageSubscriberKey, func(c di.Container) (any, error) {
		return am.NewMessageSubscriber(
			stream,
		), nil
	})

	container.AddScoped(constants.EventPublisherKey, func(c di.Container) (any, error) {
		reg := c.Get(constants.RegistryKey).(registry.Registry)
		msgPublisher := c.Get(constants.MessagePublisherKey).(am.MessagePublisher)

		return am.NewEventPublisher(reg, msgPublisher), nil
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

	container.AddScoped(constants.OrdersRepoKey, func(c di.Container) (any, error) {
		tx := c.Get(constants.DatabaseTransactionKey).(*sql.Tx)
		reg := c.Get(constants.RegistryKey).(registry.Registry)

		eventStore := pg.NewEventStore(constants.EventsTableName, tx, reg)
		snapshotStore := pg.NewSnapshotStore(constants.SnapshotsTableName, tx, reg)
		aggregateStore := es.AggregateStoreWithMiddleware(eventStore, snapshotStore)

		return es.NewAggregateRepository[*domain.Order](domain.OrderAggregate, reg, aggregateStore), nil
	})

	// setup application
	container.AddScoped(constants.ApplicationKey, func(c di.Container) (any, error) {
		orders := c.Get(constants.OrdersRepoKey).(domain.OrderRepository)
		domainDispatcher := c.Get(constants.DomainDispatcherKey).(*ddd.EventDispatcher[ddd.Event])

		app := application.New(orders, domainDispatcher)

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
		reg := c.Get(constants.RegistryKey).(registry.Registry)
		app := c.Get(constants.ApplicationKey).(application.App)

		integrationEventhandlers := handlers.NewIntegrationEventHandlers(app)
		integrationEventhandlers = logging.LogEventHandlerAccess[ddd.Event](
			integrationEventhandlers,
			"IntegrationEvents",
			service.Logger(),
		) // logging wrapper

		inboxStore := c.Get(constants.InboxStoreKey).(tm.InboxStore)
		inboxHandler := tm.InboxHandler(inboxStore)

		return am.NewEventHandler(reg, integrationEventhandlers, inboxHandler), nil
	})

	container.AddScoped(constants.CommandHandlersKey, func(c di.Container) (any, error) {
		reg := c.Get(constants.RegistryKey).(registry.Registry)
		app := c.Get(constants.ApplicationKey).(application.App)
		replyPublisher := c.Get(constants.ReplyPublisherKey).(am.ReplyPublisher)

		commandHandlers := handlers.NewCommandHandlers(app)
		commandHandlers = logging.LogCommandHandlerAccess[ddd.Command](
			commandHandlers,
			"Commands",
			service.Logger(),
		) // logging wrapper

		inboxStore := c.Get(constants.InboxStoreKey).(tm.InboxStore)
		inboxHandler := tm.InboxHandler(inboxStore)

		return am.NewCommandHandler(reg, replyPublisher, commandHandlers, inboxHandler), nil
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

	if err := handlers.RegisterIntegrationEventHandlersTx(container); err != nil {
		return err
	}
	if err := handlers.RegisterCommandHandlersTx(container); err != nil {
		return err
	}

	startOutboxProcessor(ctx, outboxProcessor, service.Logger())

	return nil
}

func startOutboxProcessor(ctx context.Context, outboxProcessor tm.OutboxProcessor, logger *slog.Logger) {
	go func() {
		err := outboxProcessor.Start(ctx)
		if err != nil {
			logger.Error("ordering outbox processor encountered an error", "error", err.Error())
		}
	}()
}
