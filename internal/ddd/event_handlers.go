package ddd

import "context"

type EventHandler[T Event] interface {
	HandleEvent(ctx context.Context, event T) error
}

type EventHandlerFunc[T Event] func(ctx context.Context, event T) error

func (f EventHandlerFunc[T]) HandleEvent(ctx context.Context, event T) error {
	return f(ctx, event)
}

type eventHandler[T Event] struct {
	h       EventHandler[T]
	filters map[string]struct{}
}
