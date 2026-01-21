package logging

import (
	"context"
	"fmt"
	"log/slog"
	"mall/internal/ddd"
)

type EventHandler[T ddd.Event] struct {
	ddd.EventHandler[T]
	label  string
	logger *slog.Logger
}

func LogEventHandlerAccess[T ddd.Event](handlers ddd.EventHandler[T], label string, logger *slog.Logger) EventHandler[T] {
	return EventHandler[T]{
		EventHandler: handlers,
		label:        label,
		logger:       logger,
	}
}

func (h EventHandler[T]) HandleEvent(ctx context.Context, event T) (err error) {
	h.logger.Info(fmt.Sprintf("--> Customers.%s.On(%s)", h.label, event.EventName()))
	defer func() {
		if err != nil {
			h.logger.Error(err.Error())
		}
		h.logger.Info(fmt.Sprintf("<-- Customers.%s.On(%s)", h.label, event.EventName()))
	}()

	return h.EventHandler.HandleEvent(ctx, event)
}
