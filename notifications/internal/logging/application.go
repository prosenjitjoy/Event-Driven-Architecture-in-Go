package logging

import (
	"context"
	"log/slog"
	"mall/notifications/internal/application"
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

func (a Application) NotifyOrderCreated(ctx context.Context, notify application.OrderCreated) (err error) {
	a.logger.InfoContext(ctx, "--> Notifications.NotifyOrderCreated")
	defer func() {
		if err != nil {
			a.logger.ErrorContext(ctx, err.Error())
		}
		a.logger.InfoContext(ctx, "<-- Notifications.NotifyOrderCreated")
	}()

	return a.App.NotifyOrderCreated(ctx, notify)
}

func (a Application) NotifyOrderCanceled(ctx context.Context, notify application.OrderCanceled) (err error) {
	a.logger.InfoContext(ctx, "--> Notifications.NoitfyOrderCanceled")
	defer func() {
		if err != nil {
			a.logger.ErrorContext(ctx, err.Error())
		}
		a.logger.InfoContext(ctx, "<-- Notifications.NoitfyOrderCanceled")
	}()

	return a.App.NotifyOrderCanceled(ctx, notify)
}

func (a Application) NotifyOrderReady(ctx context.Context, notify application.OrderReady) (err error) {
	a.logger.InfoContext(ctx, "--> Notifications.NotifyOrderReady")
	defer func() {
		if err != nil {
			a.logger.ErrorContext(ctx, err.Error())
		}
		a.logger.InfoContext(ctx, "<-- Notifications.NotifyOrderReady")
	}()

	return a.App.NotifyOrderReady(ctx, notify)
}
