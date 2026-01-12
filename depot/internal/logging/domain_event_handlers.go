package logging

import (
	"context"
	"log/slog"
	"mall/depot/internal/application"
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

func (h DomainEventHandlers) OnShoppingListCreated(ctx context.Context, event ddd.Event) (err error) {
	h.logger.Info("--> Depot.OnShoppingListCreated")
	defer func() {
		if err != nil {
			h.logger.Error(err.Error())
		}
		h.logger.Info("<-- Depot.OnShoppingListCreated")
	}()

	return h.DomainEventHandlers.OnShoppingListCreated(ctx, event)
}

func (h DomainEventHandlers) OnShoppingListCanceled(ctx context.Context, event ddd.Event) (err error) {
	h.logger.Info("--> Depot.OnShoppingListCanceled")
	defer func() {
		if err != nil {
			h.logger.Error(err.Error())
		}
		h.logger.Info("<-- Depot.OnShoppingListCanceled")
	}()

	return h.DomainEventHandlers.OnShoppingListCanceled(ctx, event)
}

func (h DomainEventHandlers) OnShoppingListAssigned(ctx context.Context, event ddd.Event) (err error) {
	h.logger.Info("--> Depot.OnShoppingListAssigned")
	defer func() {
		if err != nil {
			h.logger.Error(err.Error())
		}
		h.logger.Info("<-- Depot.OnShoppingListAssigned")
	}()

	return h.DomainEventHandlers.OnShoppingListAssigned(ctx, event)
}

func (h DomainEventHandlers) OnShoppingListCompleted(ctx context.Context, event ddd.Event) (err error) {
	h.logger.Info("--> Depot.OnShoppingListCompleted")
	defer func() {
		if err != nil {
			h.logger.Error(err.Error())
		}
		h.logger.Info("<-- Depot.OnShoppingListCompleted")
	}()

	return h.DomainEventHandlers.OnShoppingListCompleted(ctx, event)
}
