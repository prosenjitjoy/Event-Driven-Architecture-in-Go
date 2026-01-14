package queries

import (
	"context"
	"fmt"
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
	order, err := h.repo.Load(ctx, query.ID)
	if err != nil {
		return nil, fmt.Errorf("get order query: %w", err)
	}

	return order, nil
}
