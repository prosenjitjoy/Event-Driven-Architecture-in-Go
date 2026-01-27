package ordering

import (
	"context"
	"mall/baskets/basketspb"
	"mall/depot/depotpb"
	"mall/internal/am"
	"mall/internal/ddd"
	"mall/internal/es"
	"mall/internal/jetstream"
	"mall/internal/monolith"
	pg "mall/internal/postgres"
	"mall/internal/registry"
	"mall/internal/registry/serdes"
	"mall/ordering/internal/application"
	"mall/ordering/internal/domain"
	"mall/ordering/internal/grpc"
	"mall/ordering/internal/handlers"
	"mall/ordering/internal/logging"
	"mall/ordering/internal/rest"
	"mall/ordering/orderingpb"
)

type Module struct{}

func (m Module) Startup(ctx context.Context, mono monolith.Monolith) error {
	// setup driven adapters
	reg := registry.New()

	if err := registrations(reg); err != nil {
		return err
	}

	if err := basketspb.Registrations(reg); err != nil {
		return err
	}

	if err := orderingpb.Registrations(reg); err != nil {
		return err
	}

	if err := depotpb.Registrations(reg); err != nil {
		return err
	}

	stream := jetstream.NewStream(mono.Config().Nats.Stream, mono.JS(), mono.Logger())

	eventStream := am.NewEventStream(reg, stream)
	commandStream := am.NewCommandStream(reg, stream)

	domainDispatcher := ddd.NewEventDispatcher[ddd.Event]()

	aggregateStore := es.AggregateStoreWithMiddleware(
		pg.NewEventStore("ordering.events", mono.DB(), reg),
		pg.NewSnapshotStore("ordering.snapshots", mono.DB(), reg),
	)

	orders := es.NewAggregateRepository[*domain.Order](domain.OrderAggregate, reg, aggregateStore)

	conn, err := grpc.Dial(ctx, mono.Config().Rpc.Address())
	if err != nil {
		return err
	}

	shopping := grpc.NewShoppingListRepository(conn)

	// setup application
	app := logging.LogApplicationAccess(
		application.New(orders, shopping, domainDispatcher),
		mono.Logger(),
	)

	domainEventHandlers := logging.LogEventHandlerAccess[ddd.Event](
		handlers.NewDomainEventHandlers(eventStream),
		"DomainEvents",
		mono.Logger(),
	)

	integrationEventHandlers := logging.LogEventHandlerAccess[ddd.Event](
		handlers.NewIntegrationEventHandlers(app),
		"IntegrationEvents",
		mono.Logger(),
	)

	commandHandlers := logging.LogCommandHandlerAccess[ddd.Command](
		handlers.NewCommandHandlers(app),
		"Commands",
		mono.Logger(),
	)

	// setup driver adapters
	if err := grpc.RegisterServer(app, mono.RPC()); err != nil {
		return err
	}
	if err := rest.RegisterGateway(ctx, mono.Mux(), mono.Config().Rpc.Address()); err != nil {
		return err
	}
	if err := rest.RegisterSwagger(mono.Mux()); err != nil {
		return err
	}

	handlers.RegisterDomainEventHandlers(domainDispatcher, domainEventHandlers)

	if err := handlers.RegisterIntegrationEventHandlers(eventStream, integrationEventHandlers); err != nil {
		return err
	}

	if err := handlers.RegisterCommandHandlers(commandStream, commandHandlers); err != nil {
		return err
	}

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
