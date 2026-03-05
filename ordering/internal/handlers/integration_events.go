package handlers

import (
	"context"
	"mall/baskets/basketspb"
	"mall/depot/depotpb"
	"mall/internal/am"
	"mall/internal/ddd"
	"mall/ordering/internal/application"
	"mall/ordering/internal/application/commands"
	"mall/ordering/internal/domain"
	"time"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

type integrationHandlers[T ddd.Event] struct {
	app application.App
}

var _ ddd.EventHandler[ddd.Event] = (*integrationHandlers[ddd.Event])(nil)

func NewIntegrationEventHandlers(app application.App) ddd.EventHandler[ddd.Event] {
	return integrationHandlers[ddd.Event]{app: app}
}

func RegisterIntegrationEventHandlers(subscriber am.MessageSubscriber, handlers am.MessageHandler) error {
	_, err := subscriber.Subscribe(basketspb.BasketAggregateChannel, handlers, am.MessageFilters{
		basketspb.BasketCheckedOutEvent,
	}, am.GroupName("ordering-baskets"))
	if err != nil {
		return err
	}

	_, err = subscriber.Subscribe(depotpb.ShoppingListAggregateChannel, handlers, am.MessageFilters{
		depotpb.ShoppingListCompletedEvent,
	}, am.GroupName("ordering-depot"))
	if err != nil {
		return err
	}

	return nil
}

func (h integrationHandlers[T]) HandleEvent(ctx context.Context, event T) (err error) {
	span := trace.SpanFromContext(ctx)

	defer func(started time.Time) {
		if err != nil {
			span.AddEvent(
				"Encountered an error handling integration event",
				trace.WithAttributes(
					attribute.String("error.message", err.Error()),
				),
			)
			span.SetStatus(codes.Error, err.Error())
		}

		span.AddEvent(
			"Handled integration event",
			trace.WithAttributes(
				attribute.Int64("TookMS", int64(time.Since(started).Milliseconds())),
			),
		)
	}(time.Now())

	span.AddEvent(
		"Handling integration event",
		trace.WithAttributes(
			attribute.String("Event", event.EventName()),
		),
	)

	switch event.EventName() {
	case basketspb.BasketCheckedOutEvent:
		return h.onBasketCheckedOut(ctx, event)
	case depotpb.ShoppingListCompletedEvent:
		return h.onShoppingListCompleted(ctx, event)
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

func (h integrationHandlers[T]) onShoppingListCompleted(ctx context.Context, event ddd.Event) error {
	payload := event.Payload().(*depotpb.ShoppingListCompleted)

	return h.app.ReadyOrder(ctx, commands.ReadyOrderRequest{ID: payload.GetOrderId()})
}
