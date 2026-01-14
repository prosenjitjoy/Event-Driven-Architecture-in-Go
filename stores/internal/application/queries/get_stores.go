package queries

import (
	"context"
	"mall/stores/internal/domain"
)

type GetStoresRequest struct {
}

type GetStoresHandler struct {
	mall domain.MallRepository
}

func NewGetStoresHandler(mall domain.MallRepository) GetStoresHandler {
	return GetStoresHandler{
		mall: mall,
	}
}

func (h GetStoresHandler) GetStores(ctx context.Context, _ GetStoresRequest) ([]*domain.MallStore, error) {
	return h.mall.All(ctx)
}
