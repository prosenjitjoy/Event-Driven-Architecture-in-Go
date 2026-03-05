package logging

import (
	"context"
	"log/slog"
	"mall/baskets/internal/application"
	"mall/baskets/internal/domain"
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

func (a Application) StartBasket(ctx context.Context, start application.StartBasket) (err error) {
	a.logger.InfoContext(ctx, "--> Baskets.StartBasket")
	defer func() {
		if err != nil {
			a.logger.ErrorContext(ctx, err.Error())
		}
		a.logger.InfoContext(ctx, "<-- Baskets.StartBasket")
	}()

	return a.App.StartBasket(ctx, start)
}

func (a Application) CancelBasket(ctx context.Context, cancel application.CancelBasket) (err error) {
	a.logger.InfoContext(ctx, "--> Baskets.CancelBasket")
	defer func() {
		if err != nil {
			a.logger.ErrorContext(ctx, err.Error())
		}
		a.logger.InfoContext(ctx, "<-- Baskets.CancelBasket")
	}()

	return a.App.CancelBasket(ctx, cancel)
}

func (a Application) CheckoutBasket(ctx context.Context, checkout application.CheckoutBasket) (err error) {
	a.logger.InfoContext(ctx, "--> Baskets.CheckoutBasket")
	defer func() {
		if err != nil {
			a.logger.ErrorContext(ctx, err.Error())
		}
		a.logger.InfoContext(ctx, "<-- Baskets.CheckoutBasket")
	}()

	return a.App.CheckoutBasket(ctx, checkout)
}

func (a Application) AddItem(ctx context.Context, add application.AddItem) (err error) {
	a.logger.InfoContext(ctx, "--> Baskets.AddItem")
	defer func() {
		if err != nil {
			a.logger.ErrorContext(ctx, err.Error())
		}
		a.logger.InfoContext(ctx, "<-- Baskets.AddItem")
	}()

	return a.App.AddItem(ctx, add)
}

func (a Application) RemoveItem(ctx context.Context, remove application.RemoveItem) (err error) {
	a.logger.InfoContext(ctx, "--> Baskets.RemoveItem")
	defer func() {
		if err != nil {
			a.logger.ErrorContext(ctx, err.Error())
		}
		a.logger.InfoContext(ctx, "<-- Baskets.RemoveItem")
	}()

	return a.App.RemoveItem(ctx, remove)
}

func (a Application) GetBasket(ctx context.Context, get application.GetBasket) (basket *domain.Basket, err error) {
	a.logger.InfoContext(ctx, "--> Baskets.GetBasket")
	defer func() {
		if err != nil {
			a.logger.ErrorContext(ctx, err.Error())
		}
		a.logger.InfoContext(ctx, "<-- Baskets.GetBasket")
	}()

	return a.App.GetBasket(ctx, get)
}
