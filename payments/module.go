package payments

import (
	"context"
	"database/sql"
	"log/slog"
	"mall/internal/am"
	"mall/internal/ddd"
	"mall/internal/di"
	"mall/internal/jetstream"
	pg "mall/internal/postgres"
	"mall/internal/registry"
	"mall/internal/system"
	"mall/internal/tm"
	"mall/ordering/orderingpb"
	"mall/payments/internal/application"
	"mall/payments/internal/constants"
	"mall/payments/internal/domain"
	"mall/payments/internal/grpc"
	"mall/payments/internal/handlers"
	"mall/payments/internal/logging"
	"mall/payments/internal/postgres"
	"mall/payments/internal/rest"
	"mall/payments/paymentspb"
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

		if err := orderingpb.Registrations(reg); err != nil {
			return nil, err
		}
		if err := paymentspb.Registrations(reg); err != nil {
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

	container.AddScoped(constants.InvoicesRepoKey, func(c di.Container) (any, error) {
		tx := c.Get(constants.DatabaseTransactionKey).(*sql.Tx)

		return postgres.NewInvoiceRepository(constants.InvoicesTableName, tx), nil
	})

	container.AddScoped(constants.PaymentsRepoKey, func(c di.Container) (any, error) {
		tx := c.Get(constants.DatabaseTransactionKey).(*sql.Tx)

		return postgres.NewPaymentRepository(constants.PaymentsTableName, tx), nil
	})

	// setup application
	container.AddScoped(constants.ApplicationKey, func(c di.Container) (any, error) {
		invoices := c.Get(constants.InvoicesRepoKey).(domain.InvoiceRepository)
		payments := c.Get(constants.PaymentsRepoKey).(domain.PaymentRepository)
		domainDispatcher := c.Get(constants.DomainDispatcherKey).(*ddd.EventDispatcher[ddd.Event])

		app := application.New(invoices, payments, domainDispatcher)

		return logging.LogApplicationAccess(app, service.Logger()), nil
	})

	container.AddScoped(constants.DomainEventHandlersKey, func(c di.Container) (any, error) {
		eventPublisher := c.Get(constants.EventPublisherKey).(am.EventPublisher)

		domainEventHandlers := handlers.NewDomainEventHandlers(eventPublisher)

		return logging.LogEventHandlerAccess[ddd.Event](domainEventHandlers, "DomainEvents", service.Logger()), nil
	})

	container.AddScoped(constants.IntegrationEventHandlersKey, func(c di.Container) (any, error) {
		reg := c.Get(constants.RegistryKey).(registry.Registry)
		app := c.Get(constants.ApplicationKey).(application.App)

		integrationEventHandlers := handlers.NewIntegrationHandlers(app)
		integrationEventHandlers = logging.LogEventHandlerAccess[ddd.Event](
			integrationEventHandlers,
			"IntegrationEvents",
			service.Logger(),
		) // logging wrapper

		inboxStore := c.Get(constants.InboxStoreKey).(tm.InboxStore)
		inboxHandler := tm.InboxHandler(inboxStore)

		return am.NewEventHandler(reg, integrationEventHandlers, inboxHandler), nil
	})

	container.AddScoped(constants.CommandHandlersKey, func(c di.Container) (any, error) {
		reg := c.Get(constants.RegistryKey).(registry.Registry)
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
			logger.Error("payments outbox processor encountered an error", "error", err.Error())
		}
	}()
}
