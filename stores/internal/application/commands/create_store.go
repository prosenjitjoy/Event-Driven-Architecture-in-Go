package commands

import (
	"context"
	"mall/internal/ddd"
	"mall/stores/internal/domain"
)

type CreateStoreRequest struct {
	ID       string
	Name     string
	Location string
}

type CreateStoreHandler struct {
	stores    domain.StoreRepository
	publisher ddd.EventPublisher[ddd.Event]
}

func NewCreateStoreHandler(stores domain.StoreRepository, publisher ddd.EventPublisher[ddd.Event]) CreateStoreHandler {
	return CreateStoreHandler{
		stores:    stores,
		publisher: publisher,
	}
}

func (h CreateStoreHandler) CreateStore(ctx context.Context, cmd CreateStoreRequest) error {
	store, err := h.stores.Load(ctx, cmd.ID)
	if err != nil {
		return err
	}

	event, err := store.InitStore(cmd.Name, cmd.Location)
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
