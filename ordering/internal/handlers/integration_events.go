package handlers

import (
	"context"
	"mall/baskets/basketspb"
	"mall/internal/am"
	"mall/internal/ddd"
	"mall/ordering/internal/application"
	"mall/ordering/internal/application/commands"
	"mall/ordering/internal/domain"
)

type integrationHandlers[T ddd.Event] struct {
	app application.App
}

var _ ddd.EventHandler[ddd.Event] = (*integrationHandlers[ddd.Event])(nil)

func NewIntegrationEventHandlers(app application.App) ddd.EventHandler[ddd.Event] {
	return integrationHandlers[ddd.Event]{app: app}
}

func RegisterIntegrationEventHandlers(subscriber am.EventSubscriber, handlers ddd.EventHandler[ddd.Event]) error {
	eventMsgHandler := am.MessageHandlerFunc[am.IncomingEventMessage](func(ctx context.Context, eventMsg am.IncomingEventMessage) error {
		return handlers.HandleEvent(ctx, eventMsg)
	})

	err := subscriber.Subscribe(basketspb.BasketAggregateChannel, eventMsgHandler, am.MessageFilters{
		basketspb.BasketCheckedOutEvent,
	}, am.GroupName("ordering-baskets"))
	if err != nil {
		return err
	}

	return nil
}

func (h integrationHandlers[T]) HandleEvent(ctx context.Context, event T) error {
	switch event.EventName() {
	case basketspb.BasketCheckedOutEvent:
		return h.onBasketCheckedOut(ctx, event)
	}

	return nil
}

func (h integrationHandlers[T]) onBasketCheckedOut(ctx context.Context, event ddd.Event) error {
	payload := event.Payload().(*basketspb.BasketCheckedOut)

	items := make([]domain.Item, len(payload.GetItems()))

	for i, item := range payload.GetItems() {
		items[i] = domain.Item{
			ProductID:   item.GetProductId(),
			StoreID:     item.GetStoreId(),
			StoreName:   item.GetStoreName(),
			ProductName: item.GetProductName(),
			Price:       item.GetPrice(),
			Quantity:    int(item.GetQuantity()),
		}
	}

	err := h.app.CreateOrder(ctx, commands.CreateOrderRequest{
		ID:         payload.GetId(),
		CustomerID: payload.GetCustomerId(),
		PaymentID:  payload.GetPaymentId(),
		Items:      items,
	})
	if err != nil {
		return err
	}

	return nil
}
