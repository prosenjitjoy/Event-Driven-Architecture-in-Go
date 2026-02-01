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
	"mall/internal/monolith"
	pg "mall/internal/postgres"
	"mall/internal/registry"
	"mall/internal/registry/serdes"
	"mall/internal/tm"
	"mall/ordering/internal/application"
	"mall/ordering/internal/domain"
	"mall/ordering/internal/grpc"
	"mall/ordering/internal/handlers"
	"mall/ordering/internal/logging"
	"mall/ordering/internal/rest"
	"mall/ordering/orderingpb"
)

type Module struct{}

func (*Module) Startup(ctx context.Context, mono monolith.Monolith) error {
	container := di.New()

	// setup driven adapters
	container.AddSingleton("registry", func(c di.Container) (any, error) {
		reg := registry.New()

		if err := registrations(reg); err != nil {
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

	container.AddSingleton("conn", func(c di.Container) (any, error) {
		return grpc.Dial(ctx, mono.Config().Rpc.Address())
	})

	container.AddSingleton("outboxProcessor", func(c di.Container) (any, error) {
		db := c.Get("db").(*sql.DB)
		stream := c.Get("stream").(am.RawMessageStream)

		outboxStore := pg.NewOutboxStore("ordering.outbox", db)

		return tm.NewOutboxProcessor(stream, outboxStore), nil
	})

	container.AddScoped("tx", func(c di.Container) (any, error) {
		db := c.Get("db").(*sql.DB)

		return db.Begin()
	})

	container.AddScoped("txStream", func(c di.Container) (any, error) {
		tx := c.Get("tx").(*sql.Tx)
		stream := c.Get("stream").(am.RawMessageStream)

		outboxStore := pg.NewOutboxStore("ordering.outbox", tx)
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

		return am.NewCommandStream(reg, stream), nil
	})

	container.AddScoped("inboxMiddleware", func(c di.Container) (any, error) {
		tx := c.Get("tx").(*sql.Tx)

		inboxStore := pg.NewInboxStore("ordering.inbox", tx)

		return tm.NewInboxHandlerMiddleware(inboxStore), nil
	})

	container.AddScoped("aggregateStore", func(c di.Container) (any, error) {
		tx := c.Get("tx").(*sql.Tx)
		reg := c.Get("registry").(registry.Registry)

		eventStore := pg.NewEventStore("ordering.events", tx, reg)
		snapshotStore := pg.NewSnapshotStore("ordering.snapshots", tx, reg)

		return es.AggregateStoreWithMiddleware(eventStore, snapshotStore), nil
	})

	container.AddScoped("orders", func(c di.Container) (any, error) {
		reg := c.Get("registry").(registry.Registry)
		aggregateStore := c.Get("aggregateStore").(es.AggregateStore)

		return es.NewAggregateRepository[*domain.Order](domain.OrderAggregate, reg, aggregateStore), nil
	})

	// setup application
	container.AddScoped("app", func(c di.Container) (any, error) {
		logger := c.Get("logger").(*slog.Logger)
		orders := c.Get("orders").(domain.OrderRepository)
		domainDispatcher := c.Get("domainDispatcher").(*ddd.EventDispatcher[ddd.Event])

		app := application.New(orders, domainDispatcher)

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

		integrationEventhandlers := handlers.NewIntegrationEventHandlers(app)

		return logging.LogEventHandlerAccess[ddd.Event](integrationEventhandlers, "IntegrationEvents", logger), nil
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

func registrations(reg registry.Registry) error {
	serde := serdes.NewJsonSerde(reg)

	// order
	if err := serde.Register(domain.Order{}, func(v any) error {
		order := v.(*domain.Order)
		order.Aggregate = es.NewAggregate("", domain.OrderAggregate)
		return nil
	}); err != nil {
		return err
	}

	// order events
	if err := serde.Register(domain.OrderCreated{}); err != nil {
		return err
	}
	if err := serde.Register(domain.OrderRejected{}); err != nil {
		return err
	}
	if err := serde.Register(domain.OrderApproved{}); err != nil {
		return err
	}
	if err := serde.Register(domain.OrderCanceled{}); err != nil {
		return err
	}
	if err := serde.Register(domain.OrderReadied{}); err != nil {
		return err
	}
	if err := serde.Register(domain.OrderCompleted{}); err != nil {
		return err
	}

	// order snapshots
	if err := serde.RegisterKey(domain.OrderV1{}.SnapshotName(), domain.OrderV1{}); err != nil {
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
			logger.Error("ordering outbox processor encountered an error", "error", err.Error())
		}
	}()
}
