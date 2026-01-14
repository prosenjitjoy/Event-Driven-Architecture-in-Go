package commands

import (
	"context"
	"mall/ordering/internal/domain"
)

type CancelOrderRequest struct {
	ID string
}

type CancelOrderHandler struct {
	orders   domain.OrderRepository
	shopping domain.ShoppingRepository
}

func NewCancelOrderHandler(orders domain.OrderRepository, shopping domain.ShoppingRepository) CancelOrderHandler {
	return CancelOrderHandler{
		orders:   orders,
		shopping: shopping,
	}
}

func (h CancelOrderHandler) CancelOrder(ctx context.Context, cmd CancelOrderRequest) error {
	order, err := h.orders.Load(ctx, cmd.ID)
	if err != nil {
		return err
	}

	if err = order.Cancel(); err != nil {
		return err
	}

	if err = h.shopping.Cancel(ctx, order.ShoppingID); err != nil {
		return err
	}

	if err = h.orders.Save(ctx, order); err != nil {
		return err
	}

	return nil
}
