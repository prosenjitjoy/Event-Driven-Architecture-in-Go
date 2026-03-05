package logging

import (
	"context"
	"log/slog"
	"mall/depot/internal/application"
	"mall/depot/internal/application/commands"
	"mall/depot/internal/application/queries"
	"mall/depot/internal/domain"
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

func (a Application) CreateShoppingList(ctx context.Context, cmd commands.CreateShoppingListRequest) (err error) {
	a.logger.InfoContext(ctx, "--> Depot.CreateShoppingList")
	defer func() {
		if err != nil {
			a.logger.ErrorContext(ctx, err.Error())
		}
		a.logger.InfoContext(ctx, "<-- Depot.CreateShoppingList")
	}()

	return a.App.CreateShoppingList(ctx, cmd)
}

func (a Application) CancelShoppingList(ctx context.Context, cmd commands.CancelShoppingListRequest) (err error) {
	a.logger.InfoContext(ctx, "--> Depot.CancelShoppingList")
	defer func() {
		if err != nil {
			a.logger.ErrorContext(ctx, err.Error())
		}
		a.logger.InfoContext(ctx, "<-- Depot.CancelShoppingList")
	}()

	return a.App.CancelShoppingList(ctx, cmd)
}

func (a Application) AssignShoppingList(ctx context.Context, cmd commands.AssignShoppingListRequest) (err error) {
	a.logger.InfoContext(ctx, "--> Depot.AssignShoppingList")
	defer func() {
		if err != nil {
			a.logger.ErrorContext(ctx, err.Error())
		}
		a.logger.InfoContext(ctx, "<-- Depot.AssignShoppingList")
	}()

	return a.App.AssignShoppingList(ctx, cmd)
}

func (a Application) CompleteShoppingList(ctx context.Context, cmd commands.CompleteShoppingListRequest) (err error) {
	a.logger.InfoContext(ctx, "--> Depot.CompleteShoppingList")
	defer func() {
		if err != nil {
			a.logger.ErrorContext(ctx, err.Error())
		}
		a.logger.InfoContext(ctx, "<-- Depot.CompleteShoppingList")
	}()

	return a.App.CompleteShoppingList(ctx, cmd)
}

func (a Application) GetShoppingList(ctx context.Context, query queries.GetShoppingListRequest) (list *domain.ShoppingList, err error) {
	a.logger.InfoContext(ctx, "--> Depot.GetShoppingList")
	defer func() {
		if err != nil {
			a.logger.ErrorContext(ctx, err.Error())
		}
		a.logger.InfoContext(ctx, "<-- Depot.GetShoppingList")
	}()

	return a.App.GetShoppingList(ctx, query)
}
