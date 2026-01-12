package stores

import (
	"context"
	"mall/internal/ddd"
	"mall/internal/monolith"
	"mall/stores/internal/application"
	"mall/stores/internal/grpc"
	"mall/stores/internal/logging"
	"mall/stores/internal/postgres"
	"mall/stores/internal/rest"
)

type Module struct{}

func (m Module) Startup(ctx context.Context, mono monolith.Monolith) error {
	// setup driven adapters
	domainDispatcher := ddd.NewEventDispatcher()

	stores := postgres.NewStoreRepository("stores.stores", mono.DB())
	participatingStores := postgres.NewParticipatingStoreRepository("stores.stores", mono.DB())
	products := postgres.NewProductRepository("stores.products", mono.DB())

	// setup application
	app := logging.LogApplicationAccess(
		application.New(stores, participatingStores, products, domainDispatcher),
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

	return nil
}
