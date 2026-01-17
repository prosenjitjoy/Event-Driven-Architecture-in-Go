package application

import (
	"context"
	"fmt"
	"log/slog"
	"mall/internal/ddd"
	"mall/stores/storespb"
)

type StoreHandlers[T ddd.Event] struct {
	logger *slog.Logger
}

var _ ddd.EventHandler[ddd.Event] = (*StoreHandlers[ddd.Event])(nil)

func NewStoreHandlers(logger *slog.Logger) StoreHandlers[ddd.Event] {
	return StoreHandlers[ddd.Event]{
		logger: logger,
	}
}

func (h StoreHandlers[T]) HandleEvent(ctx context.Context, event T) error {
	switch event.EventName() {
	case storespb.StoreCreatedEvent:
		return h.onStoreCreated(ctx, event)
	case storespb.StoreParticipatingToggledEvent:
		return h.onStoreParticipationToggled(ctx, event)
	case storespb.StoreRebrandedEvent:
		return h.onStoreRebranded(ctx, event)
	}

	return nil
}

func (h StoreHandlers[T]) onStoreCreated(ctx context.Context, event ddd.Event) error {
	payload := event.Payload().(*storespb.StoreCreated)
	h.logger.Debug(fmt.Sprintf("ID: %s, Name: %s, Location: %s", payload.GetId(), payload.GetName(), payload.GetLocation()))
	return nil
}

func (h StoreHandlers[T]) onStoreParticipationToggled(ctx context.Context, event ddd.Event) error {
	payload := event.Payload().(*storespb.StoreParticipationToggled)
	h.logger.Debug(fmt.Sprintf("ID: %s, Participating: %t", payload.GetId(), payload.GetParticipating()))
	return nil
}

func (h StoreHandlers[T]) onStoreRebranded(ctx context.Context, event ddd.Event) error {
	payload := event.Payload().(*storespb.StoreRebranded)
	h.logger.Debug(fmt.Sprintf("ID: %s, Name: %s", payload.GetId(), payload.GetName()))
	return nil
}
