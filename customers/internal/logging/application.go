package logging

import (
	"context"
	"log/slog"
	"mall/customers/internal/application"
	"mall/customers/internal/domain"
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

func (a Application) RegisterCustomer(ctx context.Context, register application.RegisterCustomer) (err error) {
	a.logger.InfoContext(ctx, "--> Customers.RegisterCustomer")
	defer func() {
		if err != nil {
			a.logger.ErrorContext(ctx, err.Error())
		}
		a.logger.InfoContext(ctx, "<-- Customers.RegisterCustomer")
	}()

	return a.App.RegisterCustomer(ctx, register)
}

func (a Application) AuthorizeCustomer(ctx context.Context, authorize application.AuthorizeCustomer) (err error) {
	a.logger.InfoContext(ctx, "--> Customers.AuthorizeCustomer")
	defer func() {
		if err != nil {
			a.logger.ErrorContext(ctx, err.Error())
		}
		a.logger.InfoContext(ctx, "<-- Customers.AuthorizeCustomer")
	}()

	return a.App.AuthorizeCustomer(ctx, authorize)
}

func (a Application) GetCustomer(ctx context.Context, get application.GetCustomer) (customer *domain.Customer, err error) {
	a.logger.InfoContext(ctx, "--> Customers.GetCustomer")
	defer func() {
		if err != nil {
			a.logger.ErrorContext(ctx, err.Error())
		}
		a.logger.InfoContext(ctx, "<-- Customers.GetCustomer")
	}()

	return a.App.GetCustomer(ctx, get)
}

func (a Application) EnableCustomer(ctx context.Context, enable application.EnableCustomer) (err error) {
	a.logger.InfoContext(ctx, "--> Customers.EnableCustomer")
	defer func() {
		if err != nil {
			a.logger.ErrorContext(ctx, err.Error())
		}
		a.logger.InfoContext(ctx, "<-- Customers.EnableCustomer")
	}()

	return a.App.EnableCustomer(ctx, enable)
}

func (a Application) DisableCustomer(ctx context.Context, disable application.DisableCustomer) (err error) {
	a.logger.InfoContext(ctx, "--> Customers.DisableCustomer")
	defer func() {
		if err != nil {
			a.logger.ErrorContext(ctx, err.Error())
		}
		a.logger.InfoContext(ctx, "<-- Customers.DisableCustomer")
	}()

	return a.App.DisableCustomer(ctx, disable)
}
