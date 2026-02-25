package handlers

import (
	"context"
	"mall/depot/depotpb"
	"mall/depot/internal/domain"
	"mall/internal/am"
	"mall/internal/ddd"
)

type domainHandlers[T ddd.AggregateEvent] struct {
	publisher am.EventPublisher
}

var _ ddd.EventHandler[ddd.AggregateEvent] = (*domainHandlers[ddd.AggregateEvent])(nil)

func NewDomainEventHandlers(publisher am.EventPublisher) ddd.EventHandler[ddd.AggregateEvent] {
	return domainHandlers[ddd.AggregateEvent]{
		publisher: publisher,
	}
}

func RegisterDomainEventHandlers(subscriber ddd.EventSubscriber[ddd.AggregateEvent], handler ddd.EventHandler[ddd.AggregateEvent]) {
	subscriber.Subscribe(handler, domain.ShoppingListCompletedEvent)
}

func (h domainHandlers[T]) HandleEvent(ctx context.Context, event T) error {
	switch event.EventName() {
	case domain.ShoppingListCompletedEvent:
		return h.onShoppingListCompleted(ctx, event)
	}

	return nil
}

func (h domainHandlers[T]) onShoppingListCompleted(ctx context.Context, event ddd.AggregateEvent) error {
	payload := event.Payload().(*domain.ShoppingListCompleted)

	evt := ddd.NewEvent(depotpb.ShoppingListCompletedEvent, &depotpb.ShoppingListCompleted{
		Id:      event.AggregateID(),
		OrderId: payload.ShoppingList.OrderID,
	})

	return h.publisher.Publish(ctx, depotpb.ShoppingListAggregateChannel, evt)
}
