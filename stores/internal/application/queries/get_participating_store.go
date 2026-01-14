package queries

import (
	"context"
	"mall/stores/internal/domain"
)

type GetParticipatingStoreRequest struct {
}

type GetParticipatingStoreHandler struct {
	mall domain.MallRepository
}

func NewGetParticipatingStoreHandler(mall domain.MallRepository) GetParticipatingStoreHandler {
	return GetParticipatingStoreHandler{mall: mall}
}

func (h GetParticipatingStoreHandler) GetParticipatingStores(ctx context.Context, _ GetParticipatingStoreRequest) ([]*domain.MallStore, error) {
	return h.mall.AllParticipating(ctx)
}
