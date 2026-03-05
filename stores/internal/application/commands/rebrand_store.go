package commands

import (
	"context"
	"mall/internal/ddd"
	"mall/stores/internal/domain"
)

type RebrandStoreRequest struct {
	ID   string
	Name string
}

type RebrandStoreHandler struct {
	stores    domain.StoreRepository
	publisher ddd.EventPublisher[ddd.Event]
}

func NewRebrandStoreHandler(stores domain.StoreRepository, publisher ddd.EventPublisher[ddd.Event]) RebrandStoreHandler {
	return RebrandStoreHandler{
		stores:    stores,
		publisher: publisher,
	}
}

func (h RebrandStoreHandler) RebrandStore(ctx context.Context, cmd RebrandStoreRequest) error {
	store, err := h.stores.Load(ctx, cmd.ID)
	if err != nil {
		return err
	}

	event, err := store.Rebrand(cmd.Name)
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
