package queries

import (
	"context"
	"mall/stores/internal/domain"
)

type GetStoresRequest struct {
}

type GetStoresHandler struct {
	stores domain.StoreRepository
}

func NewGetStoresHandler(stores domain.StoreRepository) GetStoresHandler {
	return GetStoresHandler{
		stores: stores,
	}
}

func (h GetStoresHandler) GetStores(ctx context.Context, _ GetStoresRequest) ([]*domain.Store, error) {
	return h.stores.FindAll(ctx)
}
