package handlers

import (
	"context"
	"mall/internal/am"
	"mall/internal/ddd"
	"mall/stores/internal/domain"
	"mall/stores/storespb"
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
	subscriber.Subscribe(handlers,
		domain.StoreCreatedEvent,
		domain.StoreParticipationEnabledEvent,
		domain.StoreParticipationDisabledEvent,
		domain.StoreRebrandedEvent,
		domain.ProductAddedEvent,
		domain.ProductRebrandedEvent,
		domain.ProductPriceIncreasedEvent,
		domain.ProductPriceDecreasedEvent,
		domain.ProductRemovedEvent,
	)
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
	case domain.StoreCreatedEvent:
		return h.onStoreCreated(ctx, event)
	case domain.StoreParticipationEnabledEvent:
		return h.onStoreParticipationEnabled(ctx, event)
	case domain.StoreParticipationDisabledEvent:
		return h.onStoreParticipationDisabled(ctx, event)
	case domain.StoreRebrandedEvent:
		return h.onStoreRebranded(ctx, event)

	case domain.ProductAddedEvent:
		return h.onProductAdded(ctx, event)
	case domain.ProductRebrandedEvent:
		return h.onProductRebranded(ctx, event)
	case domain.ProductPriceIncreasedEvent:
		return h.onProductPriceIncreased(ctx, event)
	case domain.ProductPriceDecreasedEvent:
		return h.onProductPriceDecreased(ctx, event)
	case domain.ProductRemovedEvent:
		return h.onProductRemoved(ctx, event)
	}

	return nil
}

func (h domainHandlers[T]) onStoreCreated(ctx context.Context, event ddd.Event) error {
	store := event.Payload().(*domain.Store)

	evt := ddd.NewEvent(storespb.StoreCreatedEvent, &storespb.StoreCreated{
		Id:       store.ID(),
		Name:     store.Name,
		Location: store.Location,
	})

	return h.publisher.Publish(ctx, storespb.StoreAggregateChannel, evt)
}

func (h domainHandlers[T]) onStoreParticipationEnabled(ctx context.Context, event ddd.Event) error {
	store := event.Payload().(*domain.Store)

	evt := ddd.NewEvent(storespb.StoreParticipationToggledEvent, &storespb.StoreParticipationToggled{
		Id:            store.ID(),
		Participating: true,
	})

	return h.publisher.Publish(ctx, storespb.StoreAggregateChannel, evt)
}

func (h domainHandlers[T]) onStoreParticipationDisabled(ctx context.Context, event ddd.Event) error {
	store := event.Payload().(*domain.Store)

	evt := ddd.NewEvent(storespb.StoreParticipationToggledEvent, &storespb.StoreParticipationToggled{
		Id:            store.ID(),
		Participating: false,
	})

	return h.publisher.Publish(ctx, storespb.StoreAggregateChannel, evt)
}

func (h domainHandlers[T]) onStoreRebranded(ctx context.Context, event ddd.Event) error {
	store := event.Payload().(*domain.Store)

	evt := ddd.NewEvent(storespb.StoreRebrandedEvent, &storespb.StoreRebranded{
		Id:   store.ID(),
		Name: store.Name,
	})

	return h.publisher.Publish(ctx, storespb.StoreAggregateChannel, evt)
}

func (h domainHandlers[T]) onProductAdded(ctx context.Context, event ddd.Event) error {
	product := event.Payload().(*domain.Product)

	evt := ddd.NewEvent(storespb.ProductAddedEvent, &storespb.ProductAdded{
		Id:          product.ID(),
		StoreId:     product.StoreID,
		Name:        product.Name,
		Description: product.Description,
		Sku:         product.SKU,
		Price:       product.Price,
	})

	return h.publisher.Publish(ctx, storespb.ProductAggregateChannel, evt)
}

func (h domainHandlers[T]) onProductRebranded(ctx context.Context, event ddd.Event) error {
	product := event.Payload().(*domain.Product)

	evt := ddd.NewEvent(storespb.ProductRebrandedEvent, &storespb.ProductRebranded{
		Id:          product.ID(),
		Name:        product.Name,
		Description: product.Description,
	})

	return h.publisher.Publish(ctx, storespb.ProductAggregateChannel, evt)
}

func (h domainHandlers[T]) onProductPriceIncreased(ctx context.Context, event ddd.Event) error {
	payload := event.Payload().(*domain.ProductPriceDelta)

	evt := ddd.NewEvent(storespb.ProductPriceIncreasedEvent, &storespb.ProductPriceChanged{
		Id:    payload.Product.ID(),
		Delta: payload.Delta,
	})

	return h.publisher.Publish(ctx, storespb.ProductAggregateChannel, evt)
}

func (h domainHandlers[T]) onProductPriceDecreased(ctx context.Context, event ddd.Event) error {
	payload := event.Payload().(*domain.ProductPriceDelta)

	evt := ddd.NewEvent(storespb.ProductPriceDecreasedEvent, &storespb.ProductPriceChanged{
		Id:    payload.Product.ID(),
		Delta: payload.Delta,
	})

	return h.publisher.Publish(ctx, storespb.ProductAggregateChannel, evt)
}

func (h domainHandlers[T]) onProductRemoved(ctx context.Context, event ddd.Event) error {
	product := event.Payload().(*domain.Product)

	evt := ddd.NewEvent(storespb.ProductRemovedEvent, &storespb.ProductRemoved{
		Id: product.ID(),
	})

	return h.publisher.Publish(ctx, storespb.ProductAggregateChannel, evt)
}
