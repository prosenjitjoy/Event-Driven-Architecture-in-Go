package commands

import (
	"context"
	"mall/stores/internal/domain"
)

type CreateStoreRequest struct {
	ID       string
	Name     string
	Location string
}

type CreateStoreHandler struct {
	stores domain.StoreRepository
}

func NewCreateStoreHandler(stores domain.StoreRepository) CreateStoreHandler {
	return CreateStoreHandler{
		stores: stores,
	}
}

func (h CreateStoreHandler) CreateStore(ctx context.Context, cmd CreateStoreRequest) error {
	store, err := domain.CreateStore(cmd.ID, cmd.Name, cmd.Location)
	if err != nil {
		return err
	}

	return h.stores.Save(ctx, store)
}
