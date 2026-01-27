package payments

import (
	"context"
	"mall/internal/am"
	"mall/internal/ddd"
	"mall/internal/jetstream"
	"mall/internal/monolith"
	"mall/internal/registry"
	"mall/ordering/orderingpb"
	"mall/payments/internal/application"
	"mall/payments/internal/grpc"
	"mall/payments/internal/handlers"
	"mall/payments/internal/logging"
	"mall/payments/internal/postgres"
	"mall/payments/internal/rest"
	"mall/payments/paymentspb"
)

type Module struct{}

func (m Module) Startup(ctx context.Context, mono monolith.Monolith) error {
	// setup driven adapters
	reg := registry.New()

	if err := orderingpb.Registrations(reg); err != nil {
		return err
	}

	if err := paymentspb.Registrations(reg); err != nil {
		return err
	}

	stream := jetstream.NewStream(mono.Config().Nats.Stream, mono.JS(), mono.Logger())

	eventStream := am.NewEventStream(reg, stream)
	commandStream := am.NewCommandStream(reg, stream)

	domainDispatcher := ddd.NewEventDispatcher[ddd.Event]()

	invoices := postgres.NewInvoiceRepository("payments.invoices", mono.DB())
	payments := postgres.NewPaymentRepository("payments.payments", mono.DB())

	// setup application
	app := logging.LogApplicationAccess(
		application.New(invoices, payments, domainDispatcher),
		mono.Logger(),
	)

	domainEventHandlers := logging.LogEventHandlerAccess[ddd.Event](
		handlers.NewDomainEventHandlers(eventStream),
		"DomainEvents",
		mono.Logger(),
	)

	integrationEventHandlers := logging.LogEventHandlerAccess[ddd.Event](
		handlers.NewIntegrationHandlers(app),
		"IntegrationEvents",
		mono.Logger(),
	)

	commandHandlers := logging.LogCommandHandlerAccess[ddd.Command](
		handlers.NewCommandHandlers(app),
		"Commands",
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
