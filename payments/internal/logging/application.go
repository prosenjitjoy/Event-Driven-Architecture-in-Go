package logging

import (
	"context"
	"log/slog"
	"mall/payments/internal/application"
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

func (a Application) AuthorizePayment(ctx context.Context, authorize application.AuthorizePayment) (err error) {
	a.logger.Info("--> Payments.AuthorizePayment")
	defer func() {
		if err != nil {
			a.logger.Error(err.Error())
		}
		a.logger.Info("<-- Payments.AuthorizePayment")
	}()

	return a.App.AuthorizePayment(ctx, authorize)
}

func (a Application) ConfirmPayment(ctx context.Context, confirm application.ConfirmPayment) (err error) {
	a.logger.Info("--> Payments.ConfirmPayment")
	defer func() {
		if err != nil {
			a.logger.Error(err.Error())
		}
		a.logger.Info("<-- Payments.ConfirmPayment")
	}()

	return a.App.ConfirmPayment(ctx, confirm)
}

func (a Application) CreateInvoice(ctx context.Context, create application.CreateInvoice) (err error) {
	a.logger.Info("--> Payments.CreateInvoice")
	defer func() {
		if err != nil {
			a.logger.Error(err.Error())
		}
		a.logger.Info("<-- Payments.CreateInvoice")
	}()

	return a.App.CreateInvoice(ctx, create)
}

func (a Application) AdjustInvoice(ctx context.Context, adjust application.AdjustInvoice) (err error) {
	a.logger.Info("--> Payments.AdjustInvoice")
	defer func() {
		if err != nil {
			a.logger.Error(err.Error())
		}
		a.logger.Info("<-- Payments.AdjustInvoice")
	}()

	return a.App.AdjustInvoice(ctx, adjust)
}

func (a Application) PayInvoice(ctx context.Context, pay application.PayInvoice) (err error) {
	a.logger.Info("--> Payments.PayInvoice")
	defer func() {
		if err != nil {
			a.logger.Error(err.Error())
		}
		a.logger.Info("<-- Payments.PayInvoice")
	}()

	return a.App.PayInvoice(ctx, pay)
}

func (a Application) CancelInvoice(ctx context.Context, cancel application.CancelInvoice) (err error) {
	a.logger.Info("--> Payments.CancelInvoice")
	defer func() {
		if err != nil {
			a.logger.Error(err.Error())
		}
		a.logger.Info("<-- Payments.CancelInvoice")
	}()

	return a.App.CancelInvoice(ctx, cancel)
}
