package customers

import (
	"context"
	"database/sql"
	"log/slog"
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
	container.AddSingleton("registry", func(c di.Container) (any, error) {
		reg := registry.New()

		if err := customerspb.Registrations(reg); err != nil {
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

	container.AddSingleton("domainDispatcher", func(c di.Container) (any, error) {
		return ddd.NewEventDispatcher[ddd.AggregateEvent](), nil
	})

	container.AddSingleton("db", func(c di.Container) (any, error) {
		return service.DB(), nil
	})

	container.AddSingleton("outboxProcessor", func(c di.Container) (any, error) {
		db := c.Get("db").(*sql.DB)
		stream := c.Get("stream").(am.RawMessageStream)
		outboxStore := pg.NewOutboxStore("customers.outbox", db)

		return tm.NewOutboxProcessor(stream, outboxStore), nil
	})

	container.AddScoped("tx", func(c di.Container) (any, error) {
		db := c.Get("db").(*sql.DB)

		return db.Begin()
	})

	container.AddScoped("txStream", func(c di.Container) (any, error) {
		tx := c.Get("tx").(*sql.Tx)
		stream := c.Get("stream").(am.RawMessageStream)

		outboxStore := pg.NewOutboxStore("customers.outbox", tx)
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
		inboxStore := pg.NewInboxStore("customers.inbox", tx)

		return tm.NewInboxHandlerMiddleware(inboxStore), nil
	})

	container.AddScoped("customers", func(c di.Container) (any, error) {
		tx := c.Get("tx").(*sql.Tx)

		return postgres.NewCustomerRepository("customers.customers", tx), nil
	})

	// setup application
	container.AddScoped("app", func(c di.Container) (any, error) {
		logger := c.Get("logger").(*slog.Logger)
		customers := c.Get("customers").(domain.CustomerRepository)
		domainDispatcher := c.Get("domainDispatcher").(*ddd.EventDispatcher[ddd.AggregateEvent])

		app := application.New(customers, domainDispatcher)

		return logging.LogApplicationAccess(app, logger), nil
	})

	container.AddScoped("domainEventHandlers", func(c di.Container) (any, error) {
		logger := c.Get("logger").(*slog.Logger)
		eventStream := c.Get("eventStream").(am.EventStream)

		domainEventHandlers := handlers.NewDomainEventHandlers(eventStream)

		return logging.LogEventHandlerAccess[ddd.AggregateEvent](domainEventHandlers, "DomainEvents", logger), nil
	})

	container.AddScoped("commandHandlers", func(c di.Container) (any, error) {
		logger := c.Get("logger").(*slog.Logger)
		app := c.Get("app").(application.App)

		commandHandlers := handlers.NewCommandHandlers(app)

		return logging.LogCommandHandlerAccess[ddd.Command](commandHandlers, "Commands", logger), nil
	})

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

	startOutboxProcessor(ctx, container)

	return nil
}

func startOutboxProcessor(ctx context.Context, container di.Container) {
	logger := container.Get("logger").(*slog.Logger)
	outboxProcessor := container.Get("outboxProcessor").(tm.OutboxProcessor)

	go func() {
		err := outboxProcessor.Start(ctx)
		if err != nil {
			logger.Error("customers outbox processor encountered an error", "error", err.Error())
		}
	}()
}
