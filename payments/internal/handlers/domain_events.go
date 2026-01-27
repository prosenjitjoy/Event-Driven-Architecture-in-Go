package handlers

import (
	"context"
	"mall/internal/am"
	"mall/internal/ddd"
	"mall/payments/internal/domain"
	"mall/payments/paymentspb"
)

type domainHandlers[T ddd.Event] struct {
	publisher am.MessagePublisher[ddd.Event]
}

var _ ddd.EventHandler[ddd.Event] = (*domainHandlers[ddd.Event])(nil)

func NewDomainEventHandlers(publisher am.MessagePublisher[ddd.Event]) ddd.EventHandler[ddd.Event] {
	return &domainHandlers[ddd.Event]{
		publisher: publisher,
	}
}

func RegisterDomainEventHandlers(subscriber ddd.EventSubscriber[ddd.Event], handlers ddd.EventHandler[ddd.Event]) {
	subscriber.Subscribe(handlers, domain.InvoicePaidEvent)
}

func (h domainHandlers[T]) HandleEvent(ctx context.Context, event T) error {
	switch event.EventName() {
	case domain.InvoicePaidEvent:
		return h.onInvoicePaid(ctx, event)
	}

	return nil
}

func (h domainHandlers[T]) onInvoicePaid(ctx context.Context, event ddd.Event) error {
	payload := event.Payload().(*domain.InvoicePaid)

	evt := ddd.NewEvent(paymentspb.InvoicePaidEvent, &paymentspb.InvoicePaid{
		Id:      payload.ID,
		OrderId: payload.OrderID,
	})

	return h.publisher.Publish(ctx, paymentspb.InvoicePaidEvent, evt)
}
