package commands

import (
	"context"
	"mall/internal/ddd"
	"mall/ordering/internal/domain"
)

type CancelOrderRequest struct {
	ID string
}

type CancelOrderHandler struct {
	orders    domain.OrderRepository
	publisher ddd.EventPublisher[ddd.Event]
}

func NewCancelOrderHandler(orders domain.OrderRepository, publisher ddd.EventPublisher[ddd.Event]) CancelOrderHandler {
	return CancelOrderHandler{
		orders:    orders,
		publisher: publisher,
	}
}

func (h CancelOrderHandler) CancelOrder(ctx context.Context, cmd CancelOrderRequest) error {
	order, err := h.orders.Load(ctx, cmd.ID)
	if err != nil {
		return err
	}

	event, err := order.Cancel()
	if err != nil {
		return err
	}

	if err = h.orders.Save(ctx, order); err != nil {
		return err
	}

	return h.publisher.Publish(ctx, event)
}
