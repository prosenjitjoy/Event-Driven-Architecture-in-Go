package logging

import (
	"context"
	"log/slog"
	"mall/ordering/internal/application"
	"mall/ordering/internal/application/commands"
	"mall/ordering/internal/application/queries"
	"mall/ordering/internal/domain"
)

type Application struct {
	application.App
	logger *slog.Logger
}

var _ application.App = (*Application)(nil)

func LogApplicationAccess(application application.App, logger *slog.Logger) Application {
	return Application{
		App:    application,
		logger: logger,
	}
}

func (a Application) CreateOrder(ctx context.Context, cmd commands.CreateOrderRequest) (err error) {
	a.logger.InfoContext(ctx, "--> Ordering.CreateOrder")
	defer func() {
		if err != nil {
			a.logger.ErrorContext(ctx, err.Error())
		}
		a.logger.InfoContext(ctx, "<-- Ordering.CreateOrder")
	}()

	return a.App.CreateOrder(ctx, cmd)
}

func (a Application) CancelOrder(ctx context.Context, cmd commands.CancelOrderRequest) (err error) {
	a.logger.InfoContext(ctx, "--> Ordering.CancelOrder")
	defer func() {
		if err != nil {
			a.logger.ErrorContext(ctx, err.Error())
		}
		a.logger.InfoContext(ctx, "<-- Ordering.CancelOrder")
	}()

	return a.App.CancelOrder(ctx, cmd)
}

func (a Application) ReadyOrder(ctx context.Context, cmd commands.ReadyOrderRequest) (err error) {
	a.logger.InfoContext(ctx, "--> Ordering.ReadyOrder")
	defer func() {
		if err != nil {
			a.logger.ErrorContext(ctx, err.Error())
		}
		a.logger.InfoContext(ctx, "<-- Ordering.ReadyOrder")
	}()

	return a.App.ReadyOrder(ctx, cmd)
}

func (a Application) CompleteOrder(ctx context.Context, cmd commands.CompleteOrderRequest) (err error) {
	a.logger.InfoContext(ctx, "--> Ordering.CompleteOrder")
	defer func() {
		if err != nil {
			a.logger.ErrorContext(ctx, err.Error())
		}
		a.logger.InfoContext(ctx, "<-- Ordering.CompleteOrder")
	}()

	return a.App.CompleteOrder(ctx, cmd)
}

func (a Application) GetOrder(ctx context.Context, query queries.GetOrderRequest) (order *domain.Order, err error) {
	a.logger.InfoContext(ctx, "--> Ordering.GetOrder")
	defer func() {
		if err != nil {
			a.logger.ErrorContext(ctx, err.Error())
		}
		a.logger.InfoContext(ctx, "<-- Ordering.GetOrder")
	}()

	return a.App.GetOrder(ctx, query)
}
