package application

import (
	"context"
	"mall/internal/ddd"
	"mall/ordering/internal/application/commands"
	"mall/ordering/internal/application/queries"
	"mall/ordering/internal/domain"
)

type Commands interface {
	CreateOrder(ctx context.Context, cmd commands.CreateOrderRequest) error
	CancelOrder(ctx context.Context, cmd commands.CancelOrderRequest) error
	ReadyOrder(ctx context.Context, cmd commands.ReadyOrderRequest) error
	CompleteOrder(ctx context.Context, cmd commands.CompleteOrderRequest) error
}

type Queries interface {
	GetOrder(ctx context.Context, query queries.GetOrderRequest) (*domain.Order, error)
}

type App interface {
	Commands
	Queries
}

type appCommands struct {
	commands.CreateOrderHandler
	commands.CancelOrderHandler
	commands.ReadyOrderHandler
	commands.CompleteOrderHandler
}

type appQueries struct {
	queries.GetOrderHandler
}

type Application struct {
	appCommands
	appQueries
}

var _ App = (*Application)(nil)

func New(orders domain.OrderRepository, customers domain.CustomerRepository, payments domain.PaymentRepository, shopping domain.ShoppingRepository, domainPublisher ddd.EventPublisher) *Application {
	return &Application{
		appCommands: appCommands{
			CreateOrderHandler:   commands.NewCreateOrderHandler(orders, customers, payments, shopping, domainPublisher),
			CancelOrderHandler:   commands.NewCancelOrderHandler(orders, shopping, domainPublisher),
			ReadyOrderHandler:    commands.NewReadyOrderHandler(orders, domainPublisher),
			CompleteOrderHandler: commands.NewCompleteOrderHandler(orders, domainPublisher),
		},
		appQueries: appQueries{
			GetOrderHandler: queries.NewGetOrderHandler(orders),
		},
	}
}
