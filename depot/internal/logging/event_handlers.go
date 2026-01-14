package logging

import (
	"context"
	"fmt"
	"log/slog"
	"mall/internal/ddd"
)

type EventHandlers[T ddd.Event] struct {
	ddd.EventHandler[T]
	label  string
	logger *slog.Logger
}

var _ ddd.EventHandler[ddd.Event] = (*EventHandlers[ddd.Event])(nil)

func LogEventHandlerAccess[T ddd.Event](handlers ddd.EventHandler[T], label string, logger *slog.Logger) EventHandlers[T] {
	return EventHandlers[T]{
		EventHandler: handlers,
		label:        label,
		logger:       logger,
	}
}

func (h EventHandlers[T]) HandleEvent(ctx context.Context, event T) (err error) {
	h.logger.Info(fmt.Sprintf("--> Depot.%s.On(%s)", h.label, event.EventName()))
	defer func() {
		if err != nil {
			h.logger.Error(err.Error())
		}
		h.logger.Info(fmt.Sprintf("<-- Depot.%s.On(%s)", h.label, event.EventName()))
	}()

	return h.EventHandler.HandleEvent(ctx, event)
}
