package handlers

import (
	"context"
	"mall/internal/ddd"
	"mall/stores/internal/domain"
)

type mallHandlers[T ddd.AggregateEvent] struct {
	mall domain.MallRepository
}

var _ ddd.EventHandler[ddd.AggregateEvent] = (*mallHandlers[ddd.AggregateEvent])(nil)

func NewMallHandlers(mall domain.MallRepository) ddd.EventHandler[ddd.AggregateEvent] {
	return mallHandlers[ddd.AggregateEvent]{mall: mall}
}

func RegisterMallHandlers(subscriber ddd.EventSubscriber[ddd.AggregateEvent], mallHandlers ddd.EventHandler[ddd.AggregateEvent]) {
	subscriber.Subscribe(mallHandlers,
		domain.StoreCreatedEvent,
		domain.StoreParticipationEnabledEvent,
		domain.StoreParticipationDisabledEvent,
		domain.StoreRebrandedEvent,
	)
}

func (h mallHandlers[T]) HandleEvent(ctx context.Context, event T) error {
	switch event.EventName() {
	case domain.StoreCreatedEvent:
		return h.onStoreCreated(ctx, event)
	case domain.StoreParticipationEnabledEvent:
		return h.onStoreParticipationEnabled(ctx, event)
	case domain.StoreParticipationDisabledEvent:
		return h.onStoreParticipationDisabled(ctx, event)
	case domain.StoreRebrandedEvent:
		return h.onStoreRebranded(ctx, event)
	}

	return nil
}

func (h mallHandlers[T]) onStoreCreated(ctx context.Context, event ddd.AggregateEvent) error {
	payload := event.Payload().(*domain.StoreCreated)

	return h.mall.AddStore(ctx, event.AggregateID(), payload.Name, payload.Location)
}

func (h mallHandlers[T]) onStoreParticipationEnabled(ctx context.Context, event ddd.AggregateEvent) error {
	return h.mall.SetStoreParticipation(ctx, event.AggregateID(), true)
}

func (h mallHandlers[T]) onStoreParticipationDisabled(ctx context.Context, event ddd.AggregateEvent) error {
	return h.mall.SetStoreParticipation(ctx, event.AggregateID(), false)
}

func (h mallHandlers[T]) onStoreRebranded(ctx context.Context, event ddd.AggregateEvent) error {
	payload := event.Payload().(*domain.StoreRebranded)

	return h.mall.RenameStore(ctx, event.AggregateID(), payload.Name)
}
