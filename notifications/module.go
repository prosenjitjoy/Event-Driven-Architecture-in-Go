package notifications

import (
	"context"
	"mall/customers/customerspb"
	"mall/internal/am"
	"mall/internal/ddd"
	"mall/internal/jetstream"
	pg "mall/internal/postgres"
	"mall/internal/registry"
	"mall/internal/system"
	"mall/internal/tm"
	"mall/notifications/internal/application"
	"mall/notifications/internal/constants"
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

	messageSubscriber := am.NewMessageSubscriber(
		stream,
	)

	customerFallback := grpc.NewCustomerRepository(service.Config().Rpc.Service(constants.CustomersServiceName))
	customers := postgres.NewCustomerCacheRepository(
		constants.CustomersCacheTableName,
		service.DB(),
		customerFallback,
	)

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

	inboxStore := pg.NewInboxStore(constants.InboxTableName, service.DB())
	inboxHandler := tm.InboxHandler(inboxStore)

	integrationEventMsgHandlers := am.NewEventHandler(reg, integrationEventHandlers, inboxHandler)

	// setup driver adapters
	if err := grpc.RegisterServer(app, service.RPC()); err != nil {
		return err
	}

	if err := handlers.RegisterIntegrationEventHandlers(messageSubscriber, integrationEventMsgHandlers); err != nil {
		return err
	}

	return nil
}
