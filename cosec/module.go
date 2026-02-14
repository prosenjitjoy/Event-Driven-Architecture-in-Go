package cosec

import (
	"context"
	"database/sql"
	"log/slog"
	"mall/cosec/internal/application"
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
	"mall/internal/registry/serdes"
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
	container.AddSingleton("registry", func(c di.Container) (any, error) {
		reg := registry.New()

		if err := registrations(reg); err != nil {
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

	container.AddSingleton("logger", func(c di.Container) (any, error) {
		return service.Logger(), nil
	})

	container.AddSingleton("stream", func(c di.Container) (any, error) {
		logger := c.Get("logger").(*slog.Logger)

		return jetstream.NewStream(service.Config().Nats.Stream, service.JS(), logger), nil
	})

	container.AddSingleton("db", func(c di.Container) (any, error) {
		return service.DB(), nil
	})

	container.AddSingleton("outboxProcessor", func(c di.Container) (any, error) {
		db := c.Get("db").(*sql.DB)
		stream := c.Get("stream").(am.RawMessageStream)
		outboxStore := pg.NewOutboxStore("cosec.outbox", db)

		return tm.NewOutboxProcessor(stream, outboxStore), nil
	})

	container.AddScoped("tx", func(c di.Container) (any, error) {
		db := c.Get("db").(*sql.DB)

		return db.Begin()
	})

	container.AddScoped("txStream", func(c di.Container) (any, error) {
		tx := c.Get("tx").(*sql.Tx)
		outboxStore := pg.NewOutboxStore("cosec.outbox", tx)

		stream := c.Get("stream").(am.RawMessageStream)
		outboxStream := tm.NewOutboxStreamMiddleware(outboxStore)

		return am.RawMessageStreamWithMiddleware(stream, outboxStream), nil
	})

	container.AddScoped("eventStream", func(c di.Container) (any, error) {
		reg := c.Get("registry").(registry.Registry)
		stream := c.Get("txStream").(am.RawMessageStream)

		return am.NewEventStream(reg, stream), nil
	})

	container.AddScoped("commandStream", func(c di.Container) (any, error) {
		reg := c.Get("registry").(registry.Registry)
		stream := c.Get("txStream").(am.RawMessageStream)

		return am.NewCommandStream(reg, stream), nil
	})

	container.AddScoped("replyStream", func(c di.Container) (any, error) {
		reg := c.Get("registry").(registry.Registry)
		stream := c.Get("txStream").(am.RawMessageStream)

		return am.NewReplyStream(reg, stream), nil
	})

	container.AddScoped("inboxMiddleware", func(c di.Container) (any, error) {
		tx := c.Get("tx").(*sql.Tx)
		inboxStore := pg.NewInboxStore("cosec.inbox", tx)

		return tm.NewInboxHandlerMiddleware(inboxStore), nil
	})

	container.AddScoped("sagaRepo", func(c di.Container) (any, error) {
		tx := c.Get("tx").(*sql.Tx)
		reg := c.Get("registry").(registry.Registry)

		sagaStore := pg.NewSagaStore("cosec.sagas", tx, reg)

		return sec.NewSagaRepository[*domain.CreateOrderData](reg, sagaStore), nil
	})

	container.AddSingleton("saga", func(c di.Container) (any, error) {
		return application.NewCreateOrderSaga(), nil
	})

	// setup application
	container.AddScoped("orchestrator", func(c di.Container) (any, error) {
		logger := c.Get("logger").(*slog.Logger)
		saga := c.Get("saga").(sec.Saga[*domain.CreateOrderData])
		sagaRepo := c.Get("sagaRepo").(sec.SagaRepository[*domain.CreateOrderData])
		commandStream := c.Get("commandStream").(am.CommandStream)

		orchestrator := sec.NewOrchestrator[*domain.CreateOrderData](saga, sagaRepo, commandStream)

		return logging.LogReplyHandlerAccess[*domain.CreateOrderData](orchestrator, "CreateOrderSaga", logger), nil
	})

	container.AddScoped("integrationEventHandlers", func(c di.Container) (any, error) {
		logger := c.Get("logger").(*slog.Logger)
		orchestrator := c.Get("orchestrator").(sec.Orchestrator[*domain.CreateOrderData])

		integrationEventHandlers := handlers.NewIntegrationEventHandlers(orchestrator)

		return logging.LogEventHandlerAccess[ddd.Event](integrationEventHandlers, "IntegrationEvents", logger), nil
	})

	// setup driver adapters
	if err := handlers.RegisterIntegrationEventHandlersTx(container); err != nil {
		return err
	}
	if err := handlers.RegisterReplyHandlersTx(container); err != nil {
		return err
	}

	startOutboxProcessor(ctx, container)

	return nil
}

func registrations(reg registry.Registry) error {
	serde := serdes.NewJsonSerde(reg)

	// Saga data
	if err := serde.RegisterKey(application.CreateOrderSagaName, domain.CreateOrderData{}); err != nil {
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
			logger.Error("cosec outbox processor encountered an error", "error", err.Error())
		}
	}()
}
