package cosec

import (
	"context"
	"database/sql"
	"log/slog"
	"mall/cosec/internal/application"
	"mall/cosec/internal/constants"
	"mall/cosec/internal/domain"
	"mall/cosec/internal/handlers"
	"mall/cosec/internal/logging"
	"mall/customers/customerspb"
	"mall/depot/depotpb"
	"mall/internal/am"
	"mall/internal/ddd"
	"mall/internal/di"
	"mall/internal/jetstream"
	pg "mall/internal/postgres"
	"mall/internal/registry"
	"mall/internal/sec"
	"mall/internal/system"
	"mall/internal/tm"
	"mall/ordering/orderingpb"
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

		if err := domain.Registrations(reg); err != nil {
			return nil, err
		}
		if err := orderingpb.Registrations(reg); err != nil {
			return nil, err
		}
		if err := customerspb.Registrations(reg); err != nil {
			return nil, err
		}
		if err := depotpb.Registrations(reg); err != nil {
			return nil, err
		}
		if err := paymentspb.Registrations(reg); err != nil {
			return nil, err
		}

		return reg, nil
	})

	stream := jetstream.NewStream(service.Config().Nats.Stream, service.JS(), service.Logger())

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

	container.AddScoped(constants.CommandPublisherKey, func(c di.Container) (any, error) {
		reg := c.Get(constants.RegistryKey).(registry.Registry)
		msgPublisher := c.Get(constants.MessagePublisherKey).(am.MessagePublisher)

		return am.NewCommandPublisher(reg, msgPublisher), nil
	})

	container.AddScoped(constants.InboxStoreKey, func(c di.Container) (any, error) {
		tx := c.Get(constants.DatabaseTransactionKey).(*sql.Tx)

		return pg.NewInboxStore(constants.InboxTableName, tx), nil
	})

	container.AddScoped(constants.SagaStoreKey, func(c di.Container) (any, error) {
		tx := c.Get(constants.DatabaseTransactionKey).(*sql.Tx)
		reg := c.Get("registry").(registry.Registry)

		sagaStore := pg.NewSagaStore(constants.SagasTableName, tx, reg)

		return sec.NewSagaRepository[*domain.CreateOrderData](reg, sagaStore), nil
	})

	container.AddSingleton(constants.SagaKey, func(c di.Container) (any, error) {
		return application.NewCreateOrderSaga(), nil
	})

	// setup application
	container.AddScoped(constants.OrchestratorKey, func(c di.Container) (any, error) {
		saga := c.Get(constants.SagaKey).(sec.Saga[*domain.CreateOrderData])
		sagaRepo := c.Get(constants.SagaStoreKey).(sec.SagaRepository[*domain.CreateOrderData])
		commandPublisher := c.Get(constants.CommandPublisherKey).(am.CommandPublisher)

		orchestrator := sec.NewOrchestrator[*domain.CreateOrderData](saga, sagaRepo, commandPublisher)
		orchestrator = logging.LogReplyHandlerAccess[*domain.CreateOrderData](
			orchestrator,
			"CreateOrderSaga",
			service.Logger(),
		) // logging wrapper

		return orchestrator, nil
	})

	container.AddScoped(constants.IntegrationEventHandlersKey, func(c di.Container) (any, error) {
		reg := c.Get(constants.RegistryKey).(registry.Registry)
		orchestrator := c.Get(constants.OrchestratorKey).(sec.Orchestrator[*domain.CreateOrderData])

		integrationEventHandlers := handlers.NewIntegrationEventHandlers(orchestrator)
		integrationEventHandlers = logging.LogEventHandlerAccess[ddd.Event](
			integrationEventHandlers,
			"IntegrationEvents",
			service.Logger(),
		) // logging wrapper

		inboxStore := c.Get(constants.InboxStoreKey).(tm.InboxStore)
		inboxHandler := tm.InboxHandler(inboxStore)

		return am.NewEventHandler(reg, integrationEventHandlers, inboxHandler), nil
	})

	container.AddScoped(constants.ReplyHandlersKey, func(c di.Container) (any, error) {
		reg := c.Get(constants.RegistryKey).(registry.Registry)
		orchestrator := c.Get(constants.OrchestratorKey).(sec.Orchestrator[*domain.CreateOrderData])

		inboxStore := c.Get(constants.InboxStoreKey).(tm.InboxStore)
		inboxHandler := tm.InboxHandler(inboxStore)

		return handlers.NewReplyHandlers(reg, orchestrator, inboxHandler), nil
	})

	outboxStore := pg.NewOutboxStore(constants.OutboxTableName, service.DB())
	outboxProcessor := tm.NewOutboxProcessor(stream, outboxStore)

	// setup driver adapters
	if err := handlers.RegisterIntegrationEventHandlersTx(container); err != nil {
		return err
	}
	if err := handlers.RegisterReplyHandlersTx(container); err != nil {
		return err
	}

	startOutboxProcessor(ctx, outboxProcessor, service.Logger())

	return nil
}

func startOutboxProcessor(ctx context.Context, outboxProcessor tm.OutboxProcessor, logger *slog.Logger) {
	go func() {
		err := outboxProcessor.Start(ctx)
		if err != nil {
			logger.Error("cosec outbox processor encountered an error", "error", err.Error())
		}
	}()
}
