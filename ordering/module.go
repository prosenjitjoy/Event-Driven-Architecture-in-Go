package ordering

import (
	"context"
	"mall/internal/ddd"
	"mall/internal/es"
	"mall/internal/monolith"
	pg "mall/internal/postgres"
	"mall/internal/registry"
	"mall/internal/registry/serdes"
	"mall/ordering/internal/application"
	"mall/ordering/internal/domain"
	"mall/ordering/internal/grpc"
	"mall/ordering/internal/handlers"
	"mall/ordering/internal/logging"
	"mall/ordering/internal/rest"
)

type Module struct{}

func (m Module) Startup(ctx context.Context, mono monolith.Monolith) error {
	// setup driven adapters
	reg := registry.New()
	err := registrations(reg)
	if err != nil {
		return err
	}

	domainDispatcher := ddd.NewEventDispatcher[ddd.AggregateEvent]()
	aggregateStore := es.AggregateStoreWithMiddleware(
		pg.NewEventStore("ordering.events", mono.DB(), reg),
		es.NewEventPublisher(domainDispatcher),
		pg.NewSnapshotStore("ordering.snapshots", mono.DB(), reg),
	)
	orders := es.NewAggregateRepository[*domain.Order](domain.OrderAggregate, reg, aggregateStore)
	conn, err := grpc.Dial(ctx, mono.Config().Rpc.Address())
	if err != nil {
		return err
	}
	customers := grpc.NewCustomerRepository(conn)
	payments := grpc.NewPaymentRepository(conn)
	invoices := grpc.NewInvoiceRepository(conn)
	shopping := grpc.NewShoppingListRepository(conn)
	notifications := grpc.NewNotificationRepository(conn)

	// setup application
	app := logging.LogApplicationAccess(
		application.New(orders, customers, payments, shopping),
		mono.Logger(),
	)

	notificationsHandlers := logging.LogEventHandlerAccess[ddd.AggregateEvent](
		application.NewNotificationHandlers(notifications),
		"Notification",
		mono.Logger(),
	)

	invoiceHandlers := logging.LogEventHandlerAccess[ddd.AggregateEvent](
		application.NewInvoiceHandlers(invoices),
		"Invoice",
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

	handlers.RegisterNotificationHandlers(notificationsHandlers, domainDispatcher)

	handlers.RegisterInvoiceHandlers(invoiceHandlers, domainDispatcher)

	return nil
}

func registrations(reg registry.Registry) error {
	serde := serdes.NewJsonSerde(reg)

	// order
	if err := serde.Register(domain.Order{}, func(v any) error {
		order := v.(*domain.Order)
		order.Aggregate = es.NewAggregate("", domain.OrderAggregate)
		return nil
	}); err != nil {
		return err
	}

	// order events
	if err := serde.Register(domain.OrderCreated{}); err != nil {
		return err
	}
	if err := serde.Register(domain.OrderCanceled{}); err != nil {
		return err
	}
	if err := serde.Register(domain.OrderReadied{}); err != nil {
		return err
	}
	if err := serde.Register(domain.OrderCompleted{}); err != nil {
		return err
	}

	// order snapshots
	if err := serde.RegisterKey(domain.OrderV1{}.SnapshotName(), domain.OrderV1{}); err != nil {
		return err
	}

	return nil
}
