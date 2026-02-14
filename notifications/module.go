package notifications

import (
	"context"
	"mall/customers/customerspb"
	"mall/internal/am"
	"mall/internal/ddd"
	"mall/internal/jetstream"
	"mall/internal/registry"
	"mall/internal/system"
	"mall/notifications/internal/application"
	"mall/notifications/internal/grpc"
	"mall/notifications/internal/handlers"
	"mall/notifications/internal/logging"
	"mall/notifications/internal/postgres"
	"mall/ordering/orderingpb"
)

type Module struct{}

func (*Module) Startup(ctx context.Context, mono system.Service) error {
	return Root(ctx, mono)
}

func Root(ctx context.Context, service system.Service) error {
	// setup driven adapters
	reg := registry.New()

	if err := customerspb.Registrations(reg); err != nil {
		return err
	}

	if err := orderingpb.Registrations(reg); err != nil {
		return err
	}

	stream := jetstream.NewStream(service.Config().Nats.Stream, service.JS(), service.Logger())

	eventStream := am.NewEventStream(reg, stream)

	conn, err := grpc.Dial(ctx, service.Config().Rpc.Service("CUSTOMERS"))
	if err != nil {
		return err
	}

	customers := postgres.NewCustomerCacheRepository("notifications.customers_cache", service.DB(), grpc.NewCustomerRepository(conn))

	// setup application
	app := logging.LogApplicationAccess(
		application.New(customers),
		service.Logger(),
	)

	integrationEventHandlers := logging.LogEventHandlerAccess[ddd.Event](
		handlers.NewIntegrationEventHandlers(app, customers),
		"IntegrationEvents",
		service.Logger(),
	)

	// setup driver adapters
	if err := grpc.RegisterServer(app, service.RPC()); err != nil {
		return err
	}

	if err := handlers.RegisterIntegrationEventHandlers(eventStream, integrationEventHandlers); err != nil {
		return err
	}

	return nil
}
