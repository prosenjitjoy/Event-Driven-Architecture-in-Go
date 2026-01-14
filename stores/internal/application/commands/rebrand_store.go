package commands

import (
	"context"
	"mall/stores/internal/domain"
)

type RebrandStoreRequest struct {
	ID   string
	Name string
}

type RebrandStoreHandler struct {
	stores domain.StoreRepository
}

func NewRebrandStoreHandler(stores domain.StoreRepository) RebrandStoreHandler {
	return RebrandStoreHandler{
		stores: stores,
	}
}

func (h RebrandStoreHandler) RebrandStore(ctx context.Context, cmd RebrandStoreRequest) error {
	store, err := h.stores.Load(ctx, cmd.ID)
	if err != nil {
		return err
	}

	if err = store.Rebrand(cmd.Name); err != nil {
		return err
	}

	return h.stores.Save(ctx, store)
}
