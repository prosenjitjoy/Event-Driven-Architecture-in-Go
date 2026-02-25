package handlers

import (
	"context"
	"mall/customers/customerspb"
	"mall/internal/am"
	"mall/internal/ddd"
	"mall/notifications/internal/application"
	"mall/notifications/internal/domain"
	"mall/ordering/orderingpb"
)

type integrationHandlers[T ddd.Event] struct {
	app       application.App
	customers domain.CustomerCacheRepository
}

var _ ddd.EventHandler[ddd.Event] = (*integrationHandlers[ddd.Event])(nil)

func NewIntegrationEventHandlers(app application.App, customers domain.CustomerCacheRepository) ddd.EventHandler[ddd.Event] {
	return integrationHandlers[ddd.Event]{
		app:       app,
		customers: customers,
	}
}

func RegisterIntegrationEventHandlers(subscriber am.MessageSubscriber, handlers am.MessageHandler) error {
	_, err := subscriber.Subscribe(customerspb.CustomerAggregateChannel, handlers, am.MessageFilters{
		customerspb.CustomerRegisteredEvent,
		customerspb.CustomerSmsChangedEvent,
	}, am.GroupName("notification-customers"))
	if err != nil {
		return err
	}

	_, err = subscriber.Subscribe(orderingpb.OrderAggregateChannel, handlers, am.MessageFilters{
		orderingpb.OrderCreatedEvent,
		orderingpb.OrderReadiedEvent,
		orderingpb.OrderCanceledEvent,
		orderingpb.OrderCompletedEvent,
	}, am.GroupName("notification-orders"))
	if err != nil {
		return err
	}

	return nil
}

func (h integrationHandlers[T]) HandleEvent(ctx context.Context, event T) error {
	switch event.EventName() {
	case customerspb.CustomerRegisteredEvent:
		return h.onCustomerRegistered(ctx, event)
	case customerspb.CustomerSmsChangedEvent:
		return h.onCustomerSmsChanged(ctx, event)
	case orderingpb.OrderCreatedEvent:
		return h.onOrderCreated(ctx, event)
	case orderingpb.OrderReadiedEvent:
		return h.onOrderReadied(ctx, event)
	case orderingpb.OrderCanceledEvent:
		return h.onOrderCanceled(ctx, event)
	}

	return nil
}

func (h integrationHandlers[T]) onCustomerRegistered(ctx context.Context, event T) error {
	payload := event.Payload().(*customerspb.CustomerRegistered)

	return h.customers.Add(ctx, payload.GetId(), payload.GetName(), payload.GetSmsNumber())
}

func (h integrationHandlers[T]) onCustomerSmsChanged(ctx context.Context, event T) error {
	payload := event.Payload().(*customerspb.CustomerSmsChanged)

	return h.customers.UpdateSmsNumber(ctx, payload.GetId(), payload.GetSmsNumber())
}

func (h integrationHandlers[T]) onOrderCreated(ctx context.Context, event T) error {
	payload := event.Payload().(*orderingpb.OrderCreated)

	return h.app.NotifyOrderCreated(ctx, application.OrderCreated{
		OrderID:    payload.GetId(),
		CustomerID: payload.GetCustomerId(),
	})
}

func (h integrationHandlers[T]) onOrderReadied(ctx context.Context, event T) error {
	payload := event.Payload().(*orderingpb.OrderReadied)

	return h.app.NotifyOrderReady(ctx, application.OrderReady{
		OrderID:    payload.GetId(),
		CustomerID: payload.GetCustomerId(),
	})
}

func (h integrationHandlers[T]) onOrderCanceled(ctx context.Context, event T) error {
	payload := event.Payload().(*orderingpb.OrderCanceled)

	return h.app.NotifyOrderCanceled(ctx, application.OrderCanceled{
		OrderID:    payload.GetId(),
		CustomerID: payload.GetCustomerId(),
	})
}
