package commands

import (
	"context"
	"mall/ordering/internal/domain"
)

type CompleteOrderRequest struct {
	ID        string
	InvoiceID string
}

type CompleteOrderHandler struct {
	orders domain.OrderRepository
}

func NewCompleteOrderHandler(orders domain.OrderRepository) CompleteOrderHandler {
	return CompleteOrderHandler{
		orders: orders,
	}
}

func (h CompleteOrderHandler) CompleteOrder(ctx context.Context, cmd CompleteOrderRequest) error {
	order, err := h.orders.Load(ctx, cmd.ID)
	if err != nil {
		return err
	}

	if err = order.Complete(cmd.InvoiceID); err != nil {
		return err
	}

	if err = h.orders.Save(ctx, order); err != nil {
		return err
	}

	return nil
}
