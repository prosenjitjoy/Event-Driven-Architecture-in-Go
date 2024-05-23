package logging

import (
	"context"
	"log/slog"
	"mall/internal/ddd"
	"mall/ordering/internal/application"
)

type DomainEventHandlers struct {
	application.DomainEventHandlers
	logger *slog.Logger
}

var _ application.DomainEventHandlers = (*DomainEventHandlers)(nil)

func LogDomainEventHandlerAccess(handlers application.DomainEventHandlers, logger *slog.Logger) DomainEventHandlers {
	return DomainEventHandlers{
		DomainEventHandlers: handlers,
		logger:              logger,
	}
}

func (h DomainEventHandlers) OnOrderCreated(ctx context.Context, event ddd.Event) error {
	h.logger.Info("--> Ordering.OnOrderCreated")
	defer h.logger.Info("<-- Ordering.OnOrderCreated")

	return h.DomainEventHandlers.OnOrderCreated(ctx, event)
}

func (h DomainEventHandlers) OnOrderCanceled(ctx context.Context, event ddd.Event) error {
	h.logger.Info("--> Ordering.OnOrderCanceled")
	defer h.logger.Info("<-- Ordering.OnOrderCanceled")

	return h.DomainEventHandlers.OnOrderCanceled(ctx, event)
}

func (h DomainEventHandlers) OnOrderReadied(ctx context.Context, event ddd.Event) error {
	h.logger.Info("--> Ordering.OnOrderReadied")
	defer h.logger.Info("<-- Ordering.OnOrderReadied")

	return h.DomainEventHandlers.OnOrderReadied(ctx, event)
}
