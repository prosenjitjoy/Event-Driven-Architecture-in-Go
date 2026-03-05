package commands

import (
	"context"
	"mall/internal/ddd"
	"mall/stores/internal/domain"
)

type DisableParticipationRequest struct {
	ID string
}

type DisableParticipationHandler struct {
	stores    domain.StoreRepository
	publisher ddd.EventPublisher[ddd.Event]
}

func NewDisableParticipationHandler(stores domain.StoreRepository, publisher ddd.EventPublisher[ddd.Event]) DisableParticipationHandler {
	return DisableParticipationHandler{
		stores:    stores,
		publisher: publisher,
	}
}

func (h DisableParticipationHandler) DisableParticipation(ctx context.Context, cmd DisableParticipationRequest) error {
	store, err := h.stores.Load(ctx, cmd.ID)
	if err != nil {
		return err
	}

	event, err := store.DisableParticipation()
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
