package handlers

import (
	"context"
	"mall/customers/customerspb"
	"mall/customers/internal/domain"
	"mall/internal/am"
	"mall/internal/ddd"
)

type domainHandlers[T ddd.AggregateEvent] struct {
	publisher am.MessagePublisher[ddd.Event]
}

var _ ddd.EventHandler[ddd.AggregateEvent] = (*domainHandlers[ddd.AggregateEvent])(nil)

func NewDomainEventHandlers(publisher am.MessagePublisher[ddd.Event]) *domainHandlers[ddd.AggregateEvent] {
	return &domainHandlers[ddd.AggregateEvent]{
		publisher: publisher,
	}
}

func RegisterDomainEventHandlers(eventHandlers ddd.EventHandler[ddd.AggregateEvent], domainSubscriber ddd.EventSubscriber[ddd.AggregateEvent]) {
	domainSubscriber.Subscribe(eventHandlers,
		domain.CustomerRegisteredEvent,
		domain.CustomerSmsChangedEvent,
		domain.CustomerEnabledEvent,
		domain.CustomerDisabledEvent,
	)
}

func (h domainHandlers[T]) HandleEvent(ctx context.Context, event T) error {
	switch event.EventName() {
	case domain.CustomerRegisteredEvent:
		return h.onCustomerRegistered(ctx, event)
	case domain.CustomerSmsChangedEvent:
		return h.onCustomerSmsChanged(ctx, event)
	case domain.CustomerEnabledEvent:
		return h.onCustomerEnabled(ctx, event)
	case domain.CustomerDisabledEvent:
		return h.onCustomerDisabled(ctx, event)
	}

	return nil
}

func (h domainHandlers[T]) onCustomerRegistered(ctx context.Context, event ddd.AggregateEvent) error {
	payload := event.Payload().(*domain.CustomerRegistered)

	evt := ddd.NewEvent(
		customerspb.CustomerRegisteredEvent,
		&customerspb.CustomerRegistered{
			Id:        payload.Customer.ID(),
			Name:      payload.Customer.Name,
			SmsNumber: payload.Customer.SmsNumber,
		},
	)

	if err := h.publisher.Publish(ctx, customerspb.CustomerAggregateChannel, evt); err != nil {
		return err
	}

	return nil
}

func (h domainHandlers[T]) onCustomerSmsChanged(ctx context.Context, event ddd.AggregateEvent) error {
	payload := event.Payload().(*domain.CustomerSmsChanged)

	evt := ddd.NewEvent(
		customerspb.CustomerSmsChangedEvent,
		&customerspb.CustomerSmsChanged{
			Id:        payload.Customer.ID(),
			SmsNumber: payload.Customer.SmsNumber,
		},
	)

	if err := h.publisher.Publish(ctx, customerspb.CustomerAggregateChannel, evt); err != nil {
		return err
	}

	return nil
}

func (h domainHandlers[T]) onCustomerEnabled(ctx context.Context, event ddd.AggregateEvent) error {
	evt := ddd.NewEvent(
		customerspb.CustomerEnabledEvent,
		&customerspb.CustomerEnabled{Id: event.AggregateID()},
	)

	if err := h.publisher.Publish(ctx, customerspb.CustomerAggregateChannel, evt); err != nil {
		return err
	}

	return nil
}

func (h domainHandlers[T]) onCustomerDisabled(ctx context.Context, event ddd.AggregateEvent) error {
	evt := ddd.NewEvent(
		customerspb.CustomerAggregateChannel,
		&customerspb.CustomerDisabled{Id: event.AggregateID()},
	)

	if err := h.publisher.Publish(ctx, customerspb.CustomerAggregateChannel, evt); err != nil {
		return err
	}

	return nil
}
