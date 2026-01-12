package logging

import (
	"context"
	"log/slog"
	"mall/baskets/internal/application"
	"mall/internal/ddd"
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

func (h DomainEventHandlers) OnBasketStarted(ctx context.Context, event ddd.Event) (err error) {
	h.logger.Info("--> Baskets.OnBasketStarted")
	defer func() {
		if err != nil {
			h.logger.Error(err.Error())
		}
		h.logger.Info("<-- Baskets.OnBasketStarted")
	}()

	return h.DomainEventHandlers.OnBasketStarted(ctx, event)
}

func (h DomainEventHandlers) OnBasketItemAdded(ctx context.Context, event ddd.Event) (err error) {
	h.logger.Info("--> Baskets.OnBasketItemAdded")
	defer func() {
		if err != nil {
			h.logger.Error(err.Error())
		}
		h.logger.Info("<-- Baskets.OnBasketItemAdded")
	}()

	return h.DomainEventHandlers.OnBasketItemAdded(ctx, event)
}

func (h DomainEventHandlers) OnBasketItemRemoved(ctx context.Context, event ddd.Event) (err error) {
	h.logger.Info("--> Baskets.OnBasketItemRemoved")
	defer func() {
		if err != nil {
			h.logger.Error(err.Error())
		}
		h.logger.Info("<-- Baskets.OnBasketItemRemoved")
	}()

	return h.DomainEventHandlers.OnBasketItemRemoved(ctx, event)
}

func (h DomainEventHandlers) OnBasketCanceled(ctx context.Context, event ddd.Event) (err error) {
	h.logger.Info("--> Baskets.OnBasketCanceled")
	defer func() {
		if err != nil {
			h.logger.Error(err.Error())
		}
		h.logger.Info("<-- Baskets.OnBasketCanceled")
	}()

	return h.DomainEventHandlers.OnBasketCanceled(ctx, event)
}

func (h DomainEventHandlers) OnBasketCheckedOut(ctx context.Context, event ddd.Event) (err error) {
	h.logger.Info("--> Baskets.OnBasketCheckedOut")
	defer func() {
		if err != nil {
			h.logger.Error(err.Error())
		}
		h.logger.Info("<-- Baskets.OnBasketCheckedOut")
	}()

	return h.DomainEventHandlers.OnBasketCheckedOut(ctx, event)
}
