package queries

import (
	"context"
	"mall/stores/internal/domain"
)

type GetStoreRequest struct {
	ID string
}

type GetStoreHandler struct {
	stores domain.StoreRepository
}

func NewGetStoreHandler(stores domain.StoreRepository) GetStoreHandler {
	return GetStoreHandler{
		stores: stores,
	}
}

func (h GetStoreHandler) GetStore(ctx context.Context, query GetStoreRequest) (*domain.Store, error) {
	return h.stores.Find(ctx, query.ID)
}
