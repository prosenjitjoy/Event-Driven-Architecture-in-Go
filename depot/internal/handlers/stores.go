package handlers

import (
	"context"
	"mall/internal/am"
	"mall/internal/ddd"
	"mall/stores/storespb"
)

func RegisterStoreHandlers(storeHandlers ddd.EventHandler[ddd.Event], stream am.EventSubscriber) error {
	eventMsgHandler := am.MessageHandlerFunc[am.EventMessage](func(ctx context.Context, eventMsg am.EventMessage) error {
		return storeHandlers.HandleEvent(ctx, eventMsg)
	})

	return stream.Subscribe(storespb.StoreAggregateChannel, eventMsgHandler, am.MessageFilters{
		storespb.StoreCreatedEvent,
		storespb.StoreRebrandedEvent,
	}, am.GroupName("depot-stores"))
}
