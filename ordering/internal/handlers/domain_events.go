package handlers

import (
	"context"
	"mall/internal/am"
	"mall/internal/ddd"
	"mall/ordering/internal/domain"
	"mall/ordering/orderingpb"
)

type domainHandlers[T ddd.Event] struct {
	publisher am.MessagePublisher[ddd.Event]
}

func NewDomainEventHandlers(publisher am.MessagePublisher[ddd.Event]) ddd.EventHandler[ddd.Event] {
	return domainHandlers[ddd.Event]{
		publisher: publisher,
	}
}

func RegisterDomainEventHandlers(subscriber ddd.EventSubscriber[ddd.Event], handlers ddd.EventHandler[ddd.Event]) {
	subscriber.Subscribe(handlers,
		domain.OrderCreatedEvent,
		domain.OrderRejectedEvent,
		domain.OrderApprovedEvent,
		domain.OrderReadiedEvent,
		domain.OrderCanceledEvent,
		domain.OrderCompletedEvent,
	)
}

func (h domainHandlers[T]) HandleEvent(ctx context.Context, event T) error {
	switch event.EventName() {
	case domain.OrderCreatedEvent:
		return h.onOrderCreated(ctx, event)
	case domain.OrderRejectedEvent:
		return h.onOrderReadied(ctx, event)
	case domain.OrderApprovedEvent:
		return h.onOrderApproved(ctx, event)
	case domain.OrderReadiedEvent:
		return h.onOrderReadied(ctx, event)
	case domain.OrderCanceledEvent:
		return h.onOrderCanceled(ctx, event)
	case domain.OrderCompletedEvent:
		return h.onOrderCompleted(ctx, event)
	}

	return nil
}

func (h domainHandlers[T]) onOrderCreated(ctx context.Context, event ddd.Event) error {
	payload := event.Payload().(*domain.Order)

	items := make([]*orderingpb.OrderCreated_Item, len(payload.Items))

	for i, item := range payload.Items {
		items[i] = &orderingpb.OrderCreated_Item{
			ProductId: item.ProductID,
			StoreId:   item.StoreID,
			Price:     item.Price,
			Quantity:  int32(item.Quantity),
		}
	}

	evt := ddd.NewEvent(orderingpb.OrderCreatedEvent, &orderingpb.OrderCreated{
		Id:         payload.ID(),
		CustomerId: payload.CustomerID,
		PaymentId:  payload.PaymentID,
		ShoppingId: payload.ShoppingID,
		Items:      items,
	})

	return h.publisher.Publish(ctx, orderingpb.OrderAggregateChannel, evt)
}

func (h domainHandlers[T]) onOrderRejected(ctx context.Context, event ddd.Event) error {
	payload := event.Payload().(*domain.Order)

	evt := ddd.NewEvent(orderingpb.OrderRejectedEvent, &orderingpb.OrderRejected{
		Id:         payload.ID(),
		CustomerId: payload.CustomerID,
		PaymentId:  payload.PaymentID,
	})

	return h.publisher.Publish(ctx, orderingpb.OrderAggregateChannel, evt)
}

func (h domainHandlers[T]) onOrderApproved(ctx context.Context, event ddd.Event) error {
	payload := event.Payload().(*domain.Order)

	evt := ddd.NewEvent(orderingpb.OrderApprovedEvent, &orderingpb.OrderApproved{
		Id:         payload.ID(),
		CustomerId: payload.CustomerID,
		PaymentId:  payload.PaymentID,
	})

	return h.publisher.Publish(ctx, orderingpb.OrderAggregateChannel, evt)
}

func (h domainHandlers[T]) onOrderReadied(ctx context.Context, event ddd.Event) error {
	payload := event.Payload().(*domain.Order)

	evt := ddd.NewEvent(orderingpb.OrderReadiedEvent, &orderingpb.OrderReadied{
		Id:         payload.ID(),
		CustomerId: payload.CustomerID,
		PaymentId:  payload.PaymentID,
		Total:      payload.GetTotal(),
	})

	return h.publisher.Publish(ctx, orderingpb.OrderAggregateChannel, evt)
}

func (h domainHandlers[T]) onOrderCanceled(ctx context.Context, event ddd.Event) error {
	payload := event.Payload().(*domain.Order)

	evt := ddd.NewEvent(orderingpb.OrderCanceledEvent, &orderingpb.OrderCanceled{
		Id:         payload.ID(),
		CustomerId: payload.CustomerID,
		PaymentId:  payload.PaymentID,
	})

	return h.publisher.Publish(ctx, orderingpb.OrderAggregateChannel, evt)
}

func (h domainHandlers[T]) onOrderCompleted(ctx context.Context, event ddd.Event) error {
	payload := event.Payload().(*domain.Order)

	evt := ddd.NewEvent(orderingpb.OrderCompletedEvent, &orderingpb.OrderCompleted{
		Id:         payload.ID(),
		CustomerId: payload.CustomerID,
		InvoiceId:  payload.InvoiceID,
	})

	return h.publisher.Publish(ctx, orderingpb.OrderAggregateChannel, evt)
}
