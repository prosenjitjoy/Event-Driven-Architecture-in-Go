package logging

import (
	"context"
	"log/slog"
	"mall/search/internal/application"
	"mall/search/internal/domain"
)

type Application struct {
	application.Application
	logger *slog.Logger
}

var _ application.Application = (*Application)(nil)

func LogApplicationAccess(application application.Application, logger *slog.Logger) Application {
	return Application{
		Application: application,
		logger:      logger,
	}
}

func (a Application) SearchOrders(ctx context.Context, search domain.SearchOrders) (orders []*domain.Order, err error) {
	a.logger.Info("--> Search.SearchOrders")
	defer func() {
		if err != nil {
			a.logger.Error(err.Error())
		}
		a.logger.Info("<-- Search.SearchOrders")
	}()

	return a.Application.SearchOrders(ctx, search)
}

func (a Application) GetOrder(ctx context.Context, get domain.GetOrder) (order *domain.Order, err error) {
	a.logger.Info("--> Search.GetOrder")
	defer func() {
		if err != nil {
			a.logger.Error(err.Error())
		}
		a.logger.Info("<-- Search.GetOrder")
	}()

	return a.Application.GetOrder(ctx, get)
}
