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
	RejectOrder(ctx context.Context, cmd commands.RejectOrderRequest) error
	ApproveOrder(ctx context.Context, cmd commands.ApproveOrderRequest) error
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
	commands.RejectOrderHandler
	commands.ApproveOrderHandler
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

func New(orders domain.OrderRepository, publisher ddd.EventPublisher[ddd.Event]) *Application {
	return &Application{
		appCommands: appCommands{
			CreateOrderHandler:   commands.NewCreateOrderHandler(orders, publisher),
			RejectOrderHandler:   commands.NewRejectOrderHandler(orders, publisher),
			ApproveOrderHandler:  commands.NewApproveOrderHandler(orders, publisher),
			CancelOrderHandler:   commands.NewCancelOrderHandler(orders, publisher),
			ReadyOrderHandler:    commands.NewReadyOrderHandler(orders, publisher),
			CompleteOrderHandler: commands.NewCompleteOrderHandler(orders, publisher),
		},
		appQueries: appQueries{
			GetOrderHandler: queries.NewGetOrderHandler(orders),
		},
	}
}
