package handlers

import (
	"context"
	"mall/internal/am"
	"mall/internal/ddd"
	"mall/payments/internal/domain"
	"mall/payments/paymentspb"
	"time"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

type domainHandlers[T ddd.Event] struct {
	publisher am.EventPublisher
}

var _ ddd.EventHandler[ddd.Event] = (*domainHandlers[ddd.Event])(nil)

func NewDomainEventHandlers(publisher am.EventPublisher) ddd.EventHandler[ddd.Event] {
	return &domainHandlers[ddd.Event]{
		publisher: publisher,
	}
}

func RegisterDomainEventHandlers(subscriber ddd.EventSubscriber[ddd.Event], handlers ddd.EventHandler[ddd.Event]) {
	subscriber.Subscribe(handlers, domain.InvoicePaidEvent)
}

func (h domainHandlers[T]) HandleEvent(ctx context.Context, event T) (err error) {
	span := trace.SpanFromContext(ctx)

	defer func(started time.Time) {
		if err != nil {
			span.AddEvent(
				"Encountered an error handling domain event",
				trace.WithAttributes(
					attribute.String("error.message", err.Error()),
				),
			)
			span.SetStatus(codes.Error, err.Error())
		}

		span.AddEvent(
			"Handled domain event",
			trace.WithAttributes(
				attribute.Int64("TookMS", int64(time.Since(started).Milliseconds())),
			),
		)
	}(time.Now())

	span.AddEvent(
		"Handling domain event",
		trace.WithAttributes(
			attribute.String("Event", event.EventName()),
		),
	)

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
