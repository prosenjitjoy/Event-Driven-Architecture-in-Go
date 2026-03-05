package handlers

import (
	"context"
	"mall/internal/ddd"
	"mall/stores/internal/domain"
	"time"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

type CatalogHandlers[T ddd.Event] struct {
	catalog domain.CatalogRepository
}

var _ ddd.EventHandler[ddd.Event] = (*CatalogHandlers[ddd.Event])(nil)

func NewCatalogHandlers(catalog domain.CatalogRepository) ddd.EventHandler[ddd.Event] {
	return &CatalogHandlers[ddd.Event]{
		catalog: catalog,
	}
}

func RegisterCatalogHandlers(subscriber ddd.EventSubscriber[ddd.Event], catalogHandlers ddd.EventHandler[ddd.Event]) {
	subscriber.Subscribe(catalogHandlers,
		domain.ProductAddedEvent,
		domain.ProductRebrandedEvent,
		domain.ProductPriceIncreasedEvent,
		domain.ProductPriceDecreasedEvent,
		domain.ProductRemovedEvent,
	)
}

func (h CatalogHandlers[T]) HandleEvent(ctx context.Context, event T) (err error) {
	span := trace.SpanFromContext(ctx)

	defer func(started time.Time) {
		if err != nil {
			span.AddEvent(
				"Encountered an error handling catalog event",
				trace.WithAttributes(
					attribute.String("error.message", err.Error()),
				),
			)
			span.SetStatus(codes.Error, err.Error())
		}

		span.AddEvent(
			"Handled catalog event",
			trace.WithAttributes(
				attribute.Int64("TookMS", int64(time.Since(started).Milliseconds())),
			),
		)
	}(time.Now())

	span.AddEvent(
		"Handling catalog event",
		trace.WithAttributes(
			attribute.String("Event", event.EventName()),
		),
	)

	switch event.EventName() {
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

func (h CatalogHandlers[T]) onProductAdded(ctx context.Context, event ddd.Event) error {
	payload := event.Payload().(*domain.Product)

	return h.catalog.AddProduct(ctx, payload.ID(), payload.StoreID, payload.Name, payload.Description, payload.SKU, payload.Price)
}

func (h CatalogHandlers[T]) onProductRebranded(ctx context.Context, event ddd.Event) error {
	payload := event.Payload().(*domain.Product)

	return h.catalog.Rebrand(ctx, payload.ID(), payload.Name, payload.Description)
}

func (h CatalogHandlers[T]) onProductPriceIncreased(ctx context.Context, event ddd.Event) error {
	payload := event.Payload().(*domain.ProductPriceDelta)

	return h.catalog.UpdatePrice(ctx, payload.Product.ID(), payload.Delta)
}

func (h CatalogHandlers[T]) onProductPriceDecreased(ctx context.Context, event ddd.Event) error {
	payload := event.Payload().(*domain.ProductPriceDelta)

	return h.catalog.UpdatePrice(ctx, payload.Product.ID(), payload.Delta)
}

func (h CatalogHandlers[T]) onProductRemoved(ctx context.Context, event ddd.Event) error {
	payload := event.Payload().(*domain.Product)

	return h.catalog.RemoveProduct(ctx, payload.ID())
}
