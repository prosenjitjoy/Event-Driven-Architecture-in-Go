package handlers

import (
	"context"
	"mall/cosec/internal/domain"
	"mall/internal/am"
	"mall/internal/ddd"
	"mall/internal/sec"
	"mall/ordering/orderingpb"
)

type integrationHandlers[T ddd.Event] struct {
	orchestrator sec.Orchestrator[*domain.CreateOrderData]
}

var _ ddd.EventHandler[ddd.Event] = (*integrationHandlers[ddd.Event])(nil)

func NewIntegrationEventHandlers(saga sec.Orchestrator[*domain.CreateOrderData]) ddd.EventHandler[ddd.Event] {
	return integrationHandlers[ddd.Event]{
		orchestrator: saga,
	}
}

func RegisterIntegrationEventHandlers(subscriber am.MessageSubscriber, handlers am.MessageHandler) error {
	_, err := subscriber.Subscribe(orderingpb.OrderAggregateChannel, handlers, am.MessageFilters{
		orderingpb.OrderCreatedEvent,
	}, am.GroupName("cosec-ordering"))
	if err != nil {
		return err
	}

	return nil
}

func (h integrationHandlers[T]) HandleEvent(ctx context.Context, event T) error {
	switch event.EventName() {
	case orderingpb.OrderCreatedEvent:
		return h.onOrderCreated(ctx, event)
	}

	return nil
}

func (h integrationHandlers[T]) onOrderCreated(ctx context.Context, event ddd.Event) error {
	payload := event.Payload().(*orderingpb.OrderCreated)

	items := make([]domain.Item, len(payload.GetItems()))

	var total float64
	for i, item := range payload.GetItems() {
		items[i] = domain.Item{
			ProductID: item.GetProductId(),
			StoreID:   item.GetStoreId(),
			Price:     item.GetPrice(),
			Quantity:  int(item.GetQuantity()),
		}

		total += float64(item.GetQuantity()) * item.GetPrice()
	}

	data := &domain.CreateOrderData{
		OrderID:    payload.GetId(),
		CustomerID: payload.GetCustomerId(),
		PaymentID:  payload.GetPaymentId(),
		Items:      items,
		Total:      total,
	}

	// start the CreateOrderSaga
	return h.orchestrator.Start(ctx, event.ID(), data)
}
