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
	orders    domain.OrderRepository
	publisher ddd.EventPublisher[ddd.Event]
}

func NewCompleteOrderHandler(orders domain.OrderRepository, publisher ddd.EventPublisher[ddd.Event]) CompleteOrderHandler {
	return CompleteOrderHandler{
		orders:    orders,
		publisher: publisher,
	}
}

func (h CompleteOrderHandler) CompleteOrder(ctx context.Context, cmd CompleteOrderRequest) error {
	order, err := h.orders.Load(ctx, cmd.ID)
	if err != nil {
		return err
	}

	event, err := order.Complete(cmd.InvoiceID)
	if err != nil {
		return err
	}

	if err = h.orders.Save(ctx, order); err != nil {
		return err
	}

	return h.publisher.Publish(ctx, event)
}
