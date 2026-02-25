package customers

import (
	"context"
	"database/sql"
	"log/slog"
	"mall/customers/constants"
	"mall/customers/customerspb"
	"mall/customers/internal/application"
	"mall/customers/internal/domain"
	"mall/customers/internal/grpc"
	"mall/customers/internal/handlers"
	"mall/customers/internal/logging"
	"mall/customers/internal/postgres"

	"mall/customers/internal/rest"
	"mall/internal/am"
	"mall/internal/ddd"
	"mall/internal/di"
	"mall/internal/jetstream"
	pg "mall/internal/postgres"
	"mall/internal/registry"
	"mall/internal/system"
	"mall/internal/tm"
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

		if err := customerspb.Registrations(reg); err != nil {
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

	container.AddScoped(constants.CustomersRepoKey, func(c di.Container) (any, error) {
		tx := c.Get(constants.DatabaseTransactionKey).(*sql.Tx)

		return postgres.NewCustomerRepository(constants.CustomersTableName, tx), nil
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

	container.AddScoped(constants.ReplyPublisherKey, func(c di.Container) (any, error) {
		reg := c.Get(constants.RegistryKey).(registry.Registry)
		msgPublisher := c.Get(constants.MessagePublisherKey).(am.MessagePublisher)

		return am.NewReplyPublisher(reg, msgPublisher), nil
	})

	container.AddScoped(constants.InboxStoreKey, func(c di.Container) (any, error) {
		tx := c.Get(constants.DatabaseTransactionKey).(*sql.Tx)

		return pg.NewInboxStore(constants.InboxTableName, tx), nil
	})

	// setup application
	container.AddScoped(constants.ApplicationKey, func(c di.Container) (any, error) {
		customers := c.Get(constants.CustomersRepoKey).(domain.CustomerRepository)
		domainDispatcher := c.Get(constants.DomainDispatcherKey).(*ddd.EventDispatcher[ddd.AggregateEvent])

		app := application.New(customers, domainDispatcher)

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

	container.AddScoped(constants.CommandHandlersKey, func(c di.Container) (any, error) {
		reg := c.Get(constants.CommandHandlersKey).(registry.Registry)
		app := c.Get(constants.ApplicationKey).(application.App)

		commandHandlers := handlers.NewCommandHandlers(app)
		commandHandlers = logging.LogCommandHandlerAccess[ddd.Command](
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
			logger.Error("customers outbox processor encountered an error", "error", err.Error())
		}
	}()
}
