package baskets

import (
	"context"
	"mall/baskets/internal/application"
	"mall/baskets/internal/domain"
	"mall/baskets/internal/grpc"
	"mall/baskets/internal/handlers"
	"mall/baskets/internal/logging"
	"mall/baskets/internal/rest"
	"mall/internal/ddd"
	"mall/internal/es"
	"mall/internal/monolith"
	pg "mall/internal/postgres"
	"mall/internal/registry"
	"mall/internal/registry/serdes"
)

type Module struct{}

func (m *Module) Startup(ctx context.Context, mono monolith.Monolith) error {
	// setup driven adapters
	reg := registry.New()
	err := registrations(reg)
	if err != nil {
		return err
	}

	domainDispatcher := ddd.NewEventDispatcher[ddd.AggregateEvent]()
	aggregateStore := es.AggregateStoreWithMiddleware(
		pg.NewEventStore("baskets.events", mono.DB(), reg),
		es.NewEventPublisher(domainDispatcher),
		pg.NewSnapshotStore("baskets.snapshots", mono.DB(), reg),
	)

	baskets := es.NewAggregateRepository[*domain.Basket](domain.BasketAggregate, reg, aggregateStore)
	conn, err := grpc.Dial(ctx, mono.Config().Rpc.Address())
	if err != nil {
		return err
	}
	stores := grpc.NewStoreRepository(conn)
	products := grpc.NewProductRepository(conn)
	orders := grpc.NewOrderRepository(conn)

	// setup application
	app := logging.LogApplicationAccess(
		application.New(baskets, stores, products, orders),
		mono.Logger(),
	)
	orderHandlers := logging.LogEventHandlerAccess[ddd.AggregateEvent](
		application.NewOrderHandlers(orders),
		"Order",
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

	handlers.RegisterOrderHandlers(orderHandlers, domainDispatcher)

	return nil
}

func registrations(reg registry.Registry) error {
	serde := serdes.NewJsonSerde(reg)

	// Basket
	err := serde.Register(domain.Basket{}, func(v any) error {
		basket := v.(*domain.Basket)
		basket.Items = make(map[string]domain.Item)
		return nil
	})
	if err != nil {
		return err
	}

	// basket events
	if err := serde.Register(domain.BasketStarted{}); err != nil {
		return err
	}
	if err := serde.Register(domain.BasketCanceled{}); err != nil {
		return err
	}
	if err := serde.Register(domain.BasketCheckedOut{}); err != nil {
		return err
	}
	if err := serde.Register(domain.BasketItemAdded{}); err != nil {
		return err
	}
	if err := serde.Register(domain.BasketItemRemoved{}); err != nil {
		return err
	}

	// basket snapshots
	if err := serde.RegisterKey(domain.BasketV1{}.SnapshotName(), domain.BasketV1{}); err != nil {
		return err
	}

	return nil
}
