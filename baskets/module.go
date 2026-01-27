package baskets

import (
	"context"
	"mall/baskets/basketspb"
	"mall/baskets/internal/application"
	"mall/baskets/internal/domain"
	"mall/baskets/internal/grpc"
	"mall/baskets/internal/handlers"
	"mall/baskets/internal/logging"
	"mall/baskets/internal/postgres"
	"mall/baskets/internal/rest"
	"mall/internal/am"
	"mall/internal/ddd"
	"mall/internal/es"
	"mall/internal/jetstream"
	"mall/internal/monolith"
	pg "mall/internal/postgres"
	"mall/internal/registry"
	"mall/internal/registry/serdes"
	"mall/stores/storespb"
)

type Module struct{}

func (m *Module) Startup(ctx context.Context, mono monolith.Monolith) error {
	// setup driven adapters
	reg := registry.New()

	if err := registrations(reg); err != nil {
		return err
	}

	if err := basketspb.Registrations(reg); err != nil {
		return err
	}

	if err := storespb.Registrations(reg); err != nil {
		return err
	}

	stream := jetstream.NewStream(mono.Config().Nats.Stream, mono.JS(), mono.Logger())

	eventStream := am.NewEventStream(reg, stream)

	domainDispatcher := ddd.NewEventDispatcher[ddd.Event]()

	aggregateStore := es.AggregateStoreWithMiddleware(
		pg.NewEventStore("baskets.events", mono.DB(), reg),
		pg.NewSnapshotStore("baskets.snapshots", mono.DB(), reg),
	)

	baskets := es.NewAggregateRepository[*domain.Basket](domain.BasketAggregate, reg, aggregateStore)
	conn, err := grpc.Dial(ctx, mono.Config().Rpc.Address())
	if err != nil {
		return err
	}

	stores := postgres.NewStoreCacheRepository("baskets.stores_cache", mono.DB(), grpc.NewStoreRepository(conn))

	products := postgres.NewProductCacheRepository("baskets.products_cache", mono.DB(), grpc.NewProductRepository(conn))

	// setup application
	app := logging.LogApplicationAccess(
		application.New(baskets, stores, products, domainDispatcher),
		mono.Logger(),
	)

	domainEventHandlers := logging.LogEventHandlerAccess[ddd.Event](
		handlers.NewDomainEventHandlers(eventStream),
		"DomainEvents",
		mono.Logger(),
	)

	integrationEventHandlers := logging.LogEventHandlerAccess[ddd.Event](
		handlers.NewIntegrationEventHandlers(stores, products),
		"IntegrationEvents",
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

	return nil
}

func registrations(reg registry.Registry) error {
	serde := serdes.NewJsonSerde(reg)

	// Basket
	if err := serde.Register(domain.Basket{}, func(v any) error {
		basket := v.(*domain.Basket)
		basket.Items = make(map[string]domain.Item)
		return nil
	}); err != nil {
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
