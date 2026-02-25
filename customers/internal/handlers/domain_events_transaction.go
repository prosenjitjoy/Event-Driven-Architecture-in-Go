package handlers

import (
	"context"
	"mall/customers/constants"
	"mall/internal/ddd"
	"mall/internal/di"
)

func RegisterDomainEventHandlersTx(container di.Container) {
	handlers := ddd.EventHandlerFunc[ddd.AggregateEvent](func(ctx context.Context, event ddd.AggregateEvent) error {
		domainHandlers := di.Get(ctx, constants.DomainEventHandlersKey).(ddd.EventHandler[ddd.AggregateEvent])

		return domainHandlers.HandleEvent(ctx, event)
	})

	subscriber := container.Get(constants.DomainDispatcherKey).(*ddd.EventDispatcher[ddd.AggregateEvent])

	RegisterDomainEventHandlers(subscriber, handlers)
}
