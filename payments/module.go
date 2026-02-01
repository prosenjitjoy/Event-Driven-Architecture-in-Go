package payments

import (
	"context"
	"database/sql"
	"log/slog"
	"mall/internal/am"
	"mall/internal/ddd"
	"mall/internal/di"
	"mall/internal/jetstream"
	"mall/internal/monolith"
	pg "mall/internal/postgres"
	"mall/internal/registry"
	"mall/internal/tm"
	"mall/ordering/orderingpb"
	"mall/payments/internal/application"
	"mall/payments/internal/domain"
	"mall/payments/internal/grpc"
	"mall/payments/internal/handlers"
	"mall/payments/internal/logging"
	"mall/payments/internal/postgres"
	"mall/payments/internal/rest"
	"mall/payments/paymentspb"
)

type Module struct{}

func (*Module) Startup(ctx context.Context, mono monolith.Monolith) error {
	container := di.New()

	// setup driven adapters
	container.AddSingleton("registry", func(c di.Container) (any, error) {
		reg := registry.New()

		if err := orderingpb.Registrations(reg); err != nil {
			return nil, err
		}
		if err := paymentspb.Registrations(reg); err != nil {
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

	container.AddSingleton("outboxProcessor", func(c di.Container) (any, error) {
		db := c.Get("db").(*sql.DB)
		stream := c.Get("stream").(am.RawMessageStream)

		outboxStore := pg.NewOutboxStore("payments.outbox", db)

		return tm.NewOutboxProcessor(stream, outboxStore), nil
	})

	container.AddScoped("tx", func(c di.Container) (any, error) {
		db := c.Get("db").(*sql.DB)

		return db.Begin()
	})

	container.AddScoped("txStream", func(c di.Container) (any, error) {
		tx := c.Get("tx").(*sql.Tx)
		stream := c.Get("stream").(am.RawMessageStream)

		outboxStore := pg.NewOutboxStore("payments.outbox", tx)
		outboxStream := tm.NewOutboxStreamMiddleware(outboxStore)

		return am.RawMessageStreamWithMiddleware(stream, outboxStream), nil
	})

	container.AddScoped("eventStream", func(c di.Container) (any, error) {
		reg := c.Get("registry").(registry.Registry)
		stream := c.Get("txStream").(am.RawMessageStream)

		return am.NewEventStream(reg, stream), nil
	})

	container.AddScoped("replyStream", func(c di.Container) (any, error) {
		reg := c.Get("registry").(registry.Registry)
		stream := c.Get("txStream").(am.RawMessageStream)

		return am.NewReplyStream(reg, stream), nil
	})

	container.AddScoped("inboxMiddleware", func(c di.Container) (any, error) {
		tx := c.Get("tx").(*sql.Tx)

		inboxStore := pg.NewInboxStore("payments.inbox", tx)

		return tm.NewInboxHandlerMiddleware(inboxStore), nil
	})

	container.AddScoped("invoices", func(c di.Container) (any, error) {
		tx := c.Get("tx").(*sql.Tx)

		return postgres.NewInvoiceRepository("payments.invoices", tx), nil
	})

	container.AddScoped("payments", func(c di.Container) (any, error) {
		tx := c.Get("tx").(*sql.Tx)

		return postgres.NewPaymentRepository("payments.payments", tx), nil
	})

	// setup application
	container.AddScoped("app", func(c di.Container) (any, error) {
		logger := c.Get("logger").(*slog.Logger)
		invoices := c.Get("invoices").(domain.InvoiceRepository)
		payments := c.Get("payments").(domain.PaymentRepository)
		domainDispatcher := c.Get("domainDispatcher").(*ddd.EventDispatcher[ddd.Event])

		app := application.New(invoices, payments, domainDispatcher)

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
		app := c.Get("app").(application.App)

		integrationEventHandlers := handlers.NewIntegrationHandlers(app)

		return logging.LogEventHandlerAccess[ddd.Event](integrationEventHandlers, "IntegrationEvents", logger), nil
	})

	container.AddScoped("commandHandlers", func(c di.Container) (any, error) {
		logger := c.Get("logger").(*slog.Logger)
		app := c.Get("app").(application.App)

		commandHandlers := handlers.NewCommandHandlers(app)

		return logging.LogCommandHandlerAccess[ddd.Command](commandHandlers, "Commands", logger), nil
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

	if err := handlers.RegisterIntegrationEventHandlersTx(container); err != nil {
		return err
	}
	if err := handlers.RegisterCommandHandlersTx(container); err != nil {
		return err
	}

	startOutboxProcessor(ctx, container)

	return nil
}

func startOutboxProcessor(ctx context.Context, container di.Container) {
	logger := container.Get("logger").(*slog.Logger)
	outboxProcessor := container.Get("outboxProcessor").(tm.OutboxProcessor)

	go func() {
		err := outboxProcessor.Start(ctx)
		if err != nil {
			logger.Error("payments outbox processor encountered an error", "error", err.Error())
		}
	}()
}
