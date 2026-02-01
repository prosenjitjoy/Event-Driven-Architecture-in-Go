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
	"mall/ordering/orderingpb"
)

type Module struct{}

func (*Module) Startup(ctx context.Context, mono monolith.Monolith) error {

	// setup driven adapters
	reg := registry.New()

	if err := customerspb.Registrations(reg); err != nil {
		return err
	}

	if err := orderingpb.Registrations(reg); err != nil {
		return err
	}

	stream := jetstream.NewStream(mono.Config().Nats.Stream, mono.JS(), mono.Logger())

	eventStream := am.NewEventStream(reg, stream)

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

	integrationEventHandlers := logging.LogEventHandlerAccess[ddd.Event](
		handlers.NewIntegrationEventHandlers(app, customers),
		"IntegrationEvents",
		mono.Logger(),
	)

	// setup driver adapters
	if err := grpc.RegisterServer(app, mono.RPC()); err != nil {
		return err
	}

	if err := handlers.RegisterIntegrationEventHandlers(eventStream, integrationEventHandlers); err != nil {
		return err
	}

	return nil
}
