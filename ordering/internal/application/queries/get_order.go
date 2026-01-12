package queries

import (
	"context"
	"mall/ordering/internal/domain"
)

type GetOrderRequest struct {
	ID string
}

type GetOrderHandler struct {
	repo domain.OrderRepository
}

func NewGetOrderHandler(repo domain.OrderRepository) GetOrderHandler {
	return GetOrderHandler{
		repo: repo,
	}
}

func (h GetOrderHandler) GetOrder(ctx context.Context, query GetOrderRequest) (*domain.Order, error) {
	order, err := h.repo.Find(ctx, query.ID)
	if err != nil {
		return nil, err
	}

	return order, nil
}
