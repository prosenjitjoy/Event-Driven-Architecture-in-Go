package stores

import (
	"context"
	"mall/internal/ddd"
	"mall/internal/es"
	"mall/internal/monolith"
	pg "mall/internal/postgres"
	"mall/internal/registry"
	"mall/internal/registry/serdes"
	"mall/stores/internal/application"
	"mall/stores/internal/domain"
	"mall/stores/internal/grpc"
	"mall/stores/internal/handlers"
	"mall/stores/internal/logging"
	"mall/stores/internal/postgres"
	"mall/stores/internal/rest"
)

type Module struct{}

func (m Module) Startup(ctx context.Context, mono monolith.Monolith) error {
	// setup driven adapters
	reg := registry.New()
	err := registration(reg)
	if err != nil {
		return err
	}

	domainDispatcher := ddd.NewEventDispatcher[ddd.AggregateEvent]()
	aggregateStore := es.AggregateStoreWithMiddleware(
		pg.NewEventStore("stores.events", mono.DB(), reg),
		es.NewEventPublisher(domainDispatcher),
		pg.NewSnapshotStore("stores.snapshots", mono.DB(), reg),
	)
	stores := es.NewAggregateRepository[*domain.Store](domain.StoreAggregate, reg, aggregateStore)
	products := es.NewAggregateRepository[*domain.Product](domain.ProductAggregate, reg, aggregateStore)
	catalog := postgres.NewCatalogRepository("stores.products", mono.DB())
	mall := postgres.NewMallRepository("stores.stores", mono.DB())

	// setup application
	app := logging.LogApplicationAccess(
		application.New(stores, products, catalog, mall),
		mono.Logger(),
	)
	catalogHandlers := logging.LogEventHandlerAccess(
		application.NewCatalogHandlers(catalog),
		"Catalog",
		mono.Logger(),
	)
	mallHandlers := logging.LogEventHandlerAccess(
		application.NewMallHandlers(mall),
		"Mall",
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

	handlers.RegisterCatalogHandlers(catalogHandlers, domainDispatcher)
	handlers.RegisterMallHandlers(mallHandlers, domainDispatcher)

	return nil
}

func registration(reg registry.Registry) error {
	serde := serdes.NewJsonSerde(reg)

	// store
	if err := serde.Register(domain.Store{}, func(v any) error {
		store := v.(*domain.Store)
		store.Aggregate = es.NewAggregate("", domain.StoreAggregate)
		return nil
	}); err != nil {
		return err
	}

	// store events
	if err := serde.Register(domain.StoreCreated{}); err != nil {
		return err
	}
	if err := serde.RegisterKey(domain.StoreParticipationEnabledEvent, domain.StoreParticipationToggled{}); err != nil {
		return err
	}
	if err := serde.RegisterKey(domain.StoreParticipationDisabledEvent, domain.StoreParticipationToggled{}); err != nil {
		return err
	}
	if err := serde.Register(domain.StoreRebranded{}); err != nil {
		return err
	}

	// snapshots
	if err := serde.RegisterKey(domain.StoreV1{}.SnapshotName(), domain.StoreV1{}); err != nil {
		return err
	}

	// product
	if err := serde.Register(domain.Product{}, func(v any) error {
		store := v.(*domain.Product)
		store.Aggregate = es.NewAggregate("", domain.ProductAggregate)
		return nil
	}); err != nil {
		return err
	}

	// product events
	if err := serde.Register(domain.ProductAdded{}); err != nil {
		return err
	}
	if err := serde.Register(domain.ProductRebranded{}); err != nil {
		return err
	}
	if err := serde.RegisterKey(domain.ProductPriceIncreaseEvent, domain.ProductPriceChanged{}); err != nil {
		return err
	}
	if err := serde.RegisterKey(domain.ProductPriceDecreaseEvent, domain.ProductPriceChanged{}); err != nil {
		return err
	}
	if err := serde.Register(domain.ProductRemoved{}); err != nil {
		return err
	}

	// product snapshots
	if err := serde.RegisterKey(domain.ProductV1{}.SnapshotName(), domain.ProductV1{}); err != nil {
		return err
	}

	return nil
}
