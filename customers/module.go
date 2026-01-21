package customers

import (
	"context"
	"mall/customers/customerspb"
	"mall/customers/internal/application"
	"mall/customers/internal/grpc"
	"mall/customers/internal/handlers"
	"mall/customers/internal/logging"
	"mall/customers/internal/postgres"
	"mall/customers/internal/rest"
	"mall/internal/am"
	"mall/internal/ddd"
	"mall/internal/jetstream"
	"mall/internal/monolith"
	"mall/internal/registry"
)

type Module struct{}

func (m Module) Startup(ctx context.Context, mono monolith.Monolith) error {
	// setup driven adapters
	reg := registry.New()
	if err := customerspb.Registration(reg); err != nil {
		return err
	}

	eventStream := am.NewEventStream(reg, jetstream.NewStream(mono.Config().Nats.Stream, mono.JS()))

	domainDispatcher := ddd.NewEventDispatcher[ddd.AggregateEvent]()

	customers := postgres.NewCustomerRepository("customers.customers", mono.DB())

	// setup application
	app := logging.LogApplicationAccess(
		application.New(customers, domainDispatcher),
		mono.Logger(),
	)

	integrationEventHandlers := logging.LogEventHandlerAccess(
		application.NewIntegrationEventHandlers(eventStream),
		"IntegrationEvents",
		mono.Logger(),
	)

	if err := grpc.RegisterServer(app, mono.RPC()); err != nil {
		return err
	}
	if err := rest.RegisterGateway(ctx, mono.Mux(), mono.Config().Rpc.Address()); err != nil {
		return err
	}
	if err := rest.RegisterSwagger(mono.Mux()); err != nil {
		return err
	}

	handlers.RegisterIntegrationEventHandlers(integrationEventHandlers, domainDispatcher)

	return nil
}
