package commands

import (
	"context"
	"mall/internal/ddd"
	"mall/ordering/internal/domain"
)

type CompleteOrderRequest struct {
	ID        string
	InvoiceID string
}

type CompleteOrderHandler struct {
	orders          domain.OrderRepository
	domainPublisher ddd.EventPublisher
}

func NewCompleteOrderHandler(orders domain.OrderRepository, domainPublisher ddd.EventPublisher) CompleteOrderHandler {
	return CompleteOrderHandler{
		orders:          orders,
		domainPublisher: domainPublisher,
	}
}

func (h CompleteOrderHandler) CompleteOrder(ctx context.Context, cmd CompleteOrderRequest) error {
	order, err := h.orders.Find(ctx, cmd.ID)
	if err != nil {
		return err
	}

	if err = order.Complete(cmd.InvoiceID); err != nil {
		return err
	}

	if err = h.orders.Update(ctx, order); err != nil {
		return err
	}

	if err = h.domainPublisher.Publish(ctx, order.GetEvents()...); err != nil {
		return err
	}

	return nil
}
