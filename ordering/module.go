package ordering

import (
	"context"
	"mall/internal/ddd"
	"mall/internal/monolith"
	"mall/ordering/internal/application"
	"mall/ordering/internal/handlers"
	"mall/ordering/internal/logging"
	"mall/ordering/internal/postgres"
	"mall/ordering/internal/rest"
	"mall/ordering/internal/rpc"
)

type Module struct{}

func (Module) Startup(ctx context.Context, mono monolith.Monolith) error {
	// setup driven adapters
	domainDispatcher := ddd.NewEventDispatcher()
	orders := postgres.NewOrderRepository("ordering.orders", mono.DB())
	conn, err := rpc.Dial(ctx, mono.Config().Rpc.Address())
	if err != nil {
		return err
	}

	customers := rpc.NewCustomerRepository(conn)
	payments := rpc.NewPaymentRepository(conn)
	shopping := rpc.NewShoppingListRepository(conn)
	invoices := rpc.NewInvoiceRepository(conn)
	notifications := rpc.NewNotificationRepository(conn)

	// setup application
	var app application.App
	app = application.New(orders, domainDispatcher)
	app = logging.LogApplicationAccess(app, mono.Logger())

	// setup application handlers
	customersHandlers := logging.LogDomainEventHandlerAccess(application.NewCustomerRepository(customers), mono.Logger())
	paymentHandlers := logging.LogDomainEventHandlerAccess(application.NewPaymentHandlers(payments), mono.Logger())
	shoppingHandlers := logging.LogDomainEventHandlerAccess(application.NewShoppingHandlers(shopping), mono.Logger())
	invoiceHandlers := logging.LogDomainEventHandlerAccess(application.NewInvoiceHandlers(invoices), mono.Logger())
	notificationHandlers := logging.LogDomainEventHandlerAccess(application.NewNotificationHandlers(notifications), mono.Logger())

	// setup driver adapters
	if err := rpc.RegisterServer(app, mono.Rpc()); err != nil {
		return err
	}

	if err := rest.RegisterGateway(ctx, mono.Mux(), mono.Config().Rpc.Address()); err != nil {
		return err
	}

	handlers.RegisterCustomerHandlers(customersHandlers, domainDispatcher)
	handlers.RegisterPaymentHandlers(paymentHandlers, domainDispatcher)
	handlers.RegisterShoppingHandlers(shoppingHandlers, domainDispatcher)
	handlers.RegisterInvoiceHandlers(invoiceHandlers, domainDispatcher)
	handlers.RegisterNotificationHandlers(notificationHandlers, domainDispatcher)

	return nil
}
