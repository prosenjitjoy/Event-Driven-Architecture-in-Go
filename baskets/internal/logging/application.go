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
	a.logger.Info("--> Baskets.StartBasket")
	defer func() {
		if err != nil {
			a.logger.Error(err.Error())
		}
		a.logger.Info("<-- Baskets.StartBasket")
	}()

	return a.App.StartBasket(ctx, start)
}

func (a Application) CancelBasket(ctx context.Context, cancel application.CancelBasket) (err error) {
	a.logger.Info("--> Baskets.CancelBasket")
	defer func() {
		if err != nil {
			a.logger.Error(err.Error())
		}
		a.logger.Info("<-- Baskets.CancelBasket")
	}()

	return a.App.CancelBasket(ctx, cancel)
}

func (a Application) CheckoutBasket(ctx context.Context, checkout application.CheckoutBasket) (err error) {
	a.logger.Info("--> Baskets.CheckoutBasket")
	defer func() {
		if err != nil {
			a.logger.Error(err.Error())
		}
		a.logger.Info("<-- Baskets.CheckoutBasket")
	}()

	return a.App.CheckoutBasket(ctx, checkout)
}

func (a Application) AddItem(ctx context.Context, add application.AddItem) (err error) {
	a.logger.Info("--> Baskets.AddItem")
	defer func() {
		if err != nil {
			a.logger.Error(err.Error())
		}
		a.logger.Info("<-- Baskets.AddItem")
	}()

	return a.App.AddItem(ctx, add)
}

func (a Application) RemoveItem(ctx context.Context, remove application.RemoveItem) (err error) {
	a.logger.Info("--> Baskets.RemoveItem")
	defer func() {
		if err != nil {
			a.logger.Error(err.Error())
		}
		a.logger.Info("<-- Baskets.RemoveItem")
	}()

	return a.App.RemoveItem(ctx, remove)
}

func (a Application) GetBasket(ctx context.Context, get application.GetBasket) (basket *domain.Basket, err error) {
	a.logger.Info("--> Baskets.GetBasket")
	defer func() {
		if err != nil {
			a.logger.Error(err.Error())
		}
		a.logger.Info("<-- Baskets.GetBasket")
	}()

	return a.App.GetBasket(ctx, get)
}
