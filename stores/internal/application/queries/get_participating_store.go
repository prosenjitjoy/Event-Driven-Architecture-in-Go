package queries

import (
	"context"
	"mall/stores/internal/domain"
)

type GetParticipatingStoreRequest struct {
}

type GetParticipatingStoreHandler struct {
	participatingStores domain.ParticipatingStoreRepository
}

func NewGetParticipatingStoreHandler(participatingStores domain.ParticipatingStoreRepository) GetParticipatingStoreHandler {
	return GetParticipatingStoreHandler{
		participatingStores: participatingStores,
	}
}

func (h GetParticipatingStoreHandler) GetParticipatingStores(ctx context.Context, _ GetParticipatingStoreRequest) ([]*domain.Store, error) {
	return h.participatingStores.FindAll(ctx)
}
