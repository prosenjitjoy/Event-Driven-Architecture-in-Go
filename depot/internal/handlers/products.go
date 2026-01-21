package handlers

import (
	"context"
	"mall/internal/am"
	"mall/internal/ddd"
	"mall/stores/storespb"
)

func RegisterProductHandlers(productHandlers ddd.EventHandler[ddd.Event], stream am.EventSubscriber) error {
	eventMsgHandler := am.MessageHandlerFunc[am.EventMessage](func(ctx context.Context, eventMsg am.EventMessage) error {
		return productHandlers.HandleEvent(ctx, eventMsg)
	})

	return stream.Subscribe(storespb.ProductAggregateChannel, eventMsgHandler, am.MessageFilters{
		storespb.ProductAddedEvent,
		storespb.ProductRebrandedEvent,
		storespb.ProductRemovedEvent,
	}, am.GroupName("depot-products"))
}
