package handlers

import (
	"context"
	"mall/depot/depotpb"
	"mall/depot/internal/domain"
	"mall/internal/am"
	"mall/internal/ddd"
	"time"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
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
