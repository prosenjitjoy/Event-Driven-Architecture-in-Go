package commands

import (
	"context"
	"mall/stores/internal/domain"
)

type EnableParticipationRequest struct {
	ID string
}

type EnableParticipationHandler struct {
	stores domain.StoreRepository
}

func NewEnableParticipationHandler(stores domain.StoreRepository) EnableParticipationHandler {
	return EnableParticipationHandler{
		stores: stores,
	}
}

func (h EnableParticipationHandler) EnableParticipation(ctx context.Context, cmd EnableParticipationRequest) error {
	store, err := h.stores.Load(ctx, cmd.ID)
	if err != nil {
		return err
	}

	if err = store.EnableParticipation(); err != nil {
		return err
	}

	return h.stores.Save(ctx, store)
}
