package depot

import (
	"context"
	"mall/depot/depotpb"
	"mall/depot/internal/application"
	"mall/depot/internal/grpc"
	"mall/depot/internal/handlers"
	"mall/depot/internal/logging"
	"mall/depot/internal/postgres"
	"mall/depot/internal/rest"
	"mall/internal/am"
	"mall/internal/ddd"
	"mall/internal/jetstream"
	"mall/internal/monolith"
	"mall/internal/registry"
	"mall/stores/storespb"
)

type Module struct{}

func (m Module) Startup(ctx context.Context, mono monolith.Monolith) error {
	// setup driven adapters
	reg := registry.New()

	if err := storespb.Registrations(reg); err != nil {
		return err
	}

	if err := depotpb.Registrations(reg); err != nil {
		return err
	}

	stream := jetstream.NewStream(mono.Config().Nats.Stream, mono.JS(), mono.Logger())

	eventStream := am.NewEventStream(reg, stream)
	commandStream := am.NewCommandStream(reg, stream)

	domainDispatcher := ddd.NewEventDispatcher[ddd.AggregateEvent]()

	shoppingLists := postgres.NewShoppingListRepository("depot.shopping_lists", mono.DB())

	conn, err := grpc.Dial(ctx, mono.Config().Rpc.Address())
	if err != nil {
		return err
	}

	stores := postgres.NewStoreCacheRepository("depot.stores_cache", mono.DB(), grpc.NewStoreRepository(conn))

	products := postgres.NewProductCacheRepository("depot.products_cache", mono.DB(), grpc.NewProductRepository(conn))

	orders := grpc.NewOrderRepository(conn)

	// setup application
	app := logging.LogApplicationAccess(
		application.New(shoppingLists, stores, products, domainDispatcher),
		mono.Logger(),
	)

	domainEventHandlers := logging.LogEventHandlerAccess[ddd.AggregateEvent](
		handlers.NewDomainEventHandlers(orders),
		"DomainEvents",
		mono.Logger(),
	)

	integrationEventHandlers := logging.LogEventHandlerAccess[ddd.Event](
		handlers.NewIntegrationEventHandlers(stores, products),
		"IntegrationEvents",
		mono.Logger(),
	)

	commandHandlers := logging.LogCommandHandlerAccess[ddd.Command](
		handlers.NewCommandHandlers(app),
		"Command",
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
