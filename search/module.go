package search

import (
	"context"
	"mall/customers/customerspb"
	"mall/internal/am"
	"mall/internal/ddd"
	"mall/internal/jetstream"
	"mall/internal/monolith"
	"mall/internal/registry"
	"mall/ordering/orderingpb"
	"mall/search/internal/application"
	"mall/search/internal/grpc"
	"mall/search/internal/handlers"
	"mall/search/internal/logging"
	"mall/search/internal/postgres"
	"mall/search/internal/rest"
	"mall/stores/storespb"
)

type Module struct{}

func (m Module) Startup(ctx context.Context, mono monolith.Monolith) error {
	// setup driven adapters
	reg := registry.New()

	if err := orderingpb.Registrations(reg); err != nil {
		return err
	}

	if err := customerspb.Registration(reg); err != nil {
		return err
	}

	if err := storespb.Registrations(reg); err != nil {
		return err
	}

	eventStream := am.NewEventStream(reg, jetstream.NewStream(mono.Config().Nats.Stream, mono.JS()))

	conn, err := grpc.Dial(ctx, mono.Config().Rpc.Address())
	if err != nil {
		return err
	}

	customers := postgres.NewCustomerCacheRepository("search.customers_cache", mono.DB(), grpc.NewCustomerRepository(conn))

	stores := postgres.NewStoreCacheRepository("search.stores_cache", mono.DB(), grpc.NewStoreRepository(conn))

	products := postgres.NewProductCacheRepository("search.products_cache", mono.DB(), grpc.NewProductRepository(conn))

	orders := postgres.NewOrderRepository("search.orders", mono.DB())

	// setup application
	app := logging.LogApplicationAccess(
		application.New(orders),
		mono.Logger(),
	)

	orderHandlers := logging.LogEventHandlerAccess[ddd.Event](
		application.NewOrderHandlers(orders, customers, stores, products),
		"Order",
		mono.Logger(),
	)

	customerHandlers := logging.LogEventHandlerAccess[ddd.Event](
		application.NewCustomerHandlers(customers),
		"Customer",
		mono.Logger(),
	)

	storeHandlers := logging.LogEventHandlerAccess[ddd.Event](
		application.NewStoreHandlers(stores),
		"Store",
		mono.Logger(),
	)

	productHandlers := logging.LogEventHandlerAccess[ddd.Event](
		application.NewProductHandlers(products),
		"Product",
		mono.Logger(),
	)

	// setup driver adapters
	if err := grpc.RegisterServer(ctx, app, mono.RPC()); err != nil {
		return err
	}
	if err := rest.RegisterGateway(ctx, mono.Mux(), mono.Config().Rpc.Address()); err != nil {
		return err
	}
	if err := rest.RegisterSwagger(mono.Mux()); err != nil {
		return err
	}

	if err := handlers.RegisterOrderHandlers(orderHandlers, eventStream); err != nil {
		return err
	}
	if err := handlers.RegisterCustomerHandlers(customerHandlers, eventStream); err != nil {
		return err
	}
	if err := handlers.RegisterStoreHandlers(storeHandlers, eventStream); err != nil {
		return err
	}
	if err := handlers.RegisterProductHandlers(productHandlers, eventStream); err != nil {
		return err
	}

	return nil
}
