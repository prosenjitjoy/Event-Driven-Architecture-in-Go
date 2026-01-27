package commands

import (
	"context"
	"mall/internal/ddd"
	"mall/ordering/internal/domain"
)

type ReadyOrderRequest struct {
	ID string
}

type ReadyOrderHandler struct {
	orders    domain.OrderRepository
	publisher ddd.EventPublisher[ddd.Event]
}

func NewReadyOrderHandler(orders domain.OrderRepository, publisher ddd.EventPublisher[ddd.Event]) ReadyOrderHandler {
	return ReadyOrderHandler{
		orders:    orders,
		publisher: publisher,
	}
}

func (h ReadyOrderHandler) ReadyOrder(ctx context.Context, cmd ReadyOrderRequest) error {
	order, err := h.orders.Load(ctx, cmd.ID)
	if err != nil {
		return err
	}

	event, err := order.Ready()
	if err != nil {
		return err
	}

	if err = h.orders.Save(ctx, order); err != nil {
		return err
	}

	return h.publisher.Publish(ctx, event)
}
