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
)

type Module struct{}

func (m Module) Startup(ctx context.Context, mono monolith.Monolith) error {
	// setup driven adapters
	reg := registry.New()

	if err := orderingpb.Registrations(reg); err != nil {
		return err
	}

	eventStream := am.NewEventStream(reg, jetstream.NewStream(mono.Config().Nats.Stream, mono.JS()))

	domainDispatcher := ddd.NewEventDispatcher[ddd.Event]()

	invoices := postgres.NewInvoiceRepository("payments.invoices", mono.DB())
	payments := postgres.NewPaymentRepository("payments.payments", mono.DB())

	// setup application
	app := logging.LogApplicationAccess(
		application.New(invoices, payments, domainDispatcher),
		mono.Logger(),
	)

	orderHandlers := logging.LogEventHandlerAccess[ddd.Event](
		application.NewOrderHandlers(app),
		"Order",
		mono.Logger(),
	)

	integrationEventHandlers := logging.LogEventHandlerAccess[ddd.Event](
		application.NewIntegrationEventHandlers(eventStream),
		"IntegrationEvents",
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

	if err := handlers.RegisterOrderHandlers(orderHandlers, eventStream); err != nil {
		return err
	}

	handlers.RegisterIntegrationEventHandlers(integrationEventHandlers, domainDispatcher)

	return nil
}
