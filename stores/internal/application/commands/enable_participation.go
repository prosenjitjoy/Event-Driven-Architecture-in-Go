package commands

import (
	"context"
	"mall/internal/ddd"
	"mall/stores/internal/domain"
)

type EnableParticipationRequest struct {
	ID string
}

type EnableParticipationHandler struct {
	stores    domain.StoreRepository
	publisher ddd.EventPublisher[ddd.Event]
}

func NewEnableParticipationHandler(stores domain.StoreRepository, publisher ddd.EventPublisher[ddd.Event]) EnableParticipationHandler {
	return EnableParticipationHandler{
		stores:    stores,
		publisher: publisher,
	}
}

func (h EnableParticipationHandler) EnableParticipation(ctx context.Context, cmd EnableParticipationRequest) error {
	store, err := h.stores.Load(ctx, cmd.ID)
	if err != nil {
		return err
	}

	event, err := store.EnableParticipation()
	if err != nil {
		return err
	}

	if err := h.stores.Save(ctx, store); err != nil {
		return err
	}

	if err := h.publisher.Publish(ctx, event); err != nil {
		return err
	}

	return nil
}
