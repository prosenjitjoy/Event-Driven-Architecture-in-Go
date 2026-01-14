package queries

import (
	"context"
	"mall/stores/internal/domain"
)

type GetStoreRequest struct {
	ID string
}

type GetStoreHandler struct {
	mall domain.MallRepository
}

func NewGetStoreHandler(mall domain.MallRepository) GetStoreHandler {
	return GetStoreHandler{
		mall: mall,
	}
}

func (h GetStoreHandler) GetStore(ctx context.Context, query GetStoreRequest) (*domain.MallStore, error) {
	return h.mall.Find(ctx, query.ID)
}
