package application

import (
	"context"
	"mall/search/internal/domain"
)

type Application interface {
	SearchOrders(ctx context.Context, search domain.SearchOrders) ([]*domain.Order, error)
	GetOrder(ctx context.Context, get domain.GetOrder) (*domain.Order, error)
}

type app struct {
	orders domain.OrderRepository
}

var _ Application = (*app)(nil)

func New(orders domain.OrderRepository) *app {
	return &app{
		orders: orders,
	}
}

func (a app) SearchOrders(ctx context.Context, search domain.SearchOrders) ([]*domain.Order, error) {
	panic("implement me")
}

func (a app) GetOrder(ctx context.Context, get domain.GetOrder) (*domain.Order, error) {
	panic("implement me")
}
