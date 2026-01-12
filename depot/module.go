package depot

import (
	"context"
	"mall/depot/internal/application"
	"mall/depot/internal/grpc"
	"mall/depot/internal/handlers"
	"mall/depot/internal/logging"
	"mall/depot/internal/postgres"
	"mall/depot/internal/rest"
	"mall/internal/ddd"
	"mall/internal/monolith"
)

type Module struct{}

func (m Module) Startup(ctx context.Context, mono monolith.Monolith) error {
	// setup driven adapters
	domainDispatcher := ddd.NewEventDispatcher()

	shoppingLists := postgres.NewShoppingListRepository("depot.shopping_lists", mono.DB())

	conn, err := grpc.Dial(ctx, mono.Config().Rpc.Address())
	if err != nil {
		return err
	}

	stores := grpc.NewStoreRepository(conn)
	products := grpc.NewProductRepository(conn)
	orders := grpc.NewOrderRepository(conn)

	// setup application
	app := logging.LogApplicationAccess(
		application.New(shoppingLists, stores, products, domainDispatcher),
		mono.Logger(),
	)

	orderHandlers := logging.LogDomainEventHandlerAccess(
		application.NewOrderHandler(orders),
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
