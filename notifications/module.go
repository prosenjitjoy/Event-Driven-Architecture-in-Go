package notifications

import (
	"context"
	"mall/customers/customerspb"
	"mall/internal/am"
	"mall/internal/ddd"
	"mall/internal/jetstream"
	"mall/internal/monolith"
	"mall/internal/registry"
	"mall/notifications/internal/application"
	"mall/notifications/internal/grpc"
	"mall/notifications/internal/handlers"
	"mall/notifications/internal/logging"
	"mall/notifications/internal/postgres"
)

type Module struct{}

func (m Module) Startup(ctx context.Context, mono monolith.Monolith) error {
	// setup driven adapters
	reg := registry.New()
	if err := customerspb.Registration(reg); err != nil {
		return err
	}

	eventStream := am.NewEventStream(reg, jetstream.NewStream(mono.Config().Nats.Stream, mono.JS()))

	conn, err := grpc.Dial(ctx, mono.Config().Rpc.Address())
	if err != nil {
		return err
	}

	customers := postgres.NewCustomerCacheRepository("notifications.customers_cache", mono.DB(), grpc.NewCustomerRepository(conn))

	// setup application
	app := logging.LogApplicationAccess(
		application.New(customers),
		mono.Logger(),
	)

	customerHandlers := logging.LogEventHandlerAccess[ddd.Event](
		application.NewCustomerHandlers(customers),
		"Customer",
		mono.Logger(),
	)

	orderHandlers := logging.LogEventHandlerAccess[ddd.Event](
		application.NewOrderHandlers(app),
		"Order",
		mono.Logger(),
	)

	// setup driver adapters
	if err := grpc.RegisterServer(app, mono.RPC()); err != nil {
		return err
	}
	if err := handlers.RegisterCustomerHandlers(customerHandlers, eventStream); err != nil {
		return err
	}
	if err := handlers.RegisterOrderHandlers(orderHandlers, eventStream); err != nil {
		return err
	}

	return nil
}
