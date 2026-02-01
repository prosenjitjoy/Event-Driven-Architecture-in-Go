package handlers

import (
	"context"
	"mall/internal/ddd"
	"mall/internal/di"
)

func RegisterCatalogHandlersTx(container di.Container) {
	handlers := ddd.EventHandlerFunc[ddd.AggregateEvent](func(ctx context.Context, event ddd.AggregateEvent) error {
		catalogHandlers := di.Get(ctx, "catalogHandlers").(ddd.EventHandler[ddd.AggregateEvent])

		return catalogHandlers.HandleEvent(ctx, event)
	})

	subscriber := container.Get("domainDispatcher").(*ddd.EventDispatcher[ddd.AggregateEvent])

	RegisterCatalogHandlers(subscriber, handlers)
}
