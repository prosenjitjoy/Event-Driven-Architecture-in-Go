package handlers

import (
	"context"
	"mall/internal/ddd"
	"mall/internal/di"
)

func RegisterMallHandlersTx(container di.Container) {
	handlers := ddd.EventHandlerFunc[ddd.AggregateEvent](func(ctx context.Context, event ddd.AggregateEvent) error {
		mallHandlers := di.Get(ctx, "mallHandlers").(ddd.EventHandler[ddd.AggregateEvent])

		return mallHandlers.HandleEvent(ctx, event)
	})

	subscriber := container.Get("domainDispatcher").(*ddd.EventDispatcher[ddd.AggregateEvent])

	RegisterMallHandlers(subscriber, handlers)
}
