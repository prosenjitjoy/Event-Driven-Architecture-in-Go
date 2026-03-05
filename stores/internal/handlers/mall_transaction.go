package handlers

import (
	"context"
	"mall/internal/ddd"
	"mall/internal/di"
	"mall/stores/internal/constants"
)

func RegisterMallHandlersTx(container di.Container) {
	handlers := ddd.EventHandlerFunc[ddd.Event](func(ctx context.Context, event ddd.Event) error {
		mallHandlers := di.Get(ctx, constants.MallHandlersKey).(ddd.EventHandler[ddd.Event])

		return mallHandlers.HandleEvent(ctx, event)
	})

	subscriber := container.Get(constants.DomainDispatcherKey).(*ddd.EventDispatcher[ddd.Event])

	RegisterMallHandlers(subscriber, handlers)
}
