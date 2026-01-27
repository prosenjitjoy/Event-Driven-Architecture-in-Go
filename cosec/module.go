package cosec

import (
	"context"
	"mall/cosec/internal/application"
	"mall/cosec/internal/domain"
	"mall/cosec/internal/handlers"
	"mall/cosec/internal/logging"
	"mall/customers/customerspb"
	"mall/depot/depotpb"
	"mall/internal/am"
	"mall/internal/jetstream"
	"mall/internal/monolith"
	pg "mall/internal/postgres"
	"mall/internal/registry"
	"mall/internal/registry/serdes"
	"mall/internal/sec"
	"mall/ordering/orderingpb"
	"mall/payments/paymentspb"
)

type Module struct{}

func (*Module) Startup(ctx context.Context, mono monolith.Monolith) error {
	// setup driven adapters
	reg := registry.New()

	if err := registrations(reg); err != nil {
		return err
	}

	if err := orderingpb.Registrations(reg); err != nil {
		return err
	}

	if err := customerspb.Registrations(reg); err != nil {
		return err
	}

	if err := depotpb.Registrations(reg); err != nil {
		return err
	}

	if err := paymentspb.Registrations(reg); err != nil {
		return err
	}

	stream := jetstream.NewStream(mono.Config().Nats.Stream, mono.JS(), mono.Logger())

	eventStream := am.NewEventStream(reg, stream)
	commandStream := am.NewCommandStream(reg, stream)
	replyStream := am.NewReplyStream(reg, stream)

	sagaStore := pg.NewSagaStore("cosec.sagas", mono.DB(), reg)

	sagaRepo := sec.NewSagaRepository[*domain.CreateOrderData](reg, sagaStore)

	// setup application
	createOrderSEC := sec.NewOrchestrator[*domain.CreateOrderData](
		application.NewCreateOrderSaga(),
		sagaRepo,
		commandStream,
	)

	orchestrator := logging.LogReplyHandlerAccess(
		createOrderSEC,
		"CreateOrderSaga",
		mono.Logger(),
	)

	integrationEventHandlers := logging.LogEventHandlerAccess(
		handlers.NewIntegrationEventHandlers(orchestrator),
		"IntegrationEvents",
		mono.Logger(),
	)

	// setup driver adapters
	if err := handlers.RegisterIntegrationEventHandlers(eventStream, integrationEventHandlers); err != nil {
		return err
	}

	if err := handlers.RegisterReplyHandlers(replyStream, orchestrator); err != nil {
		return err
	}

	return nil
}

func registrations(reg registry.Registry) error {
	serde := serdes.NewJsonSerde(reg)

	if err := serde.RegisterKey(application.CreateOrderSagaName, domain.CreateOrderData{}); err != nil {
		return err
	}

	return nil
}
