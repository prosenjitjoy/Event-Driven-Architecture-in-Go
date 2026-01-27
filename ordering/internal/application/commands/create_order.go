package commands

import (
	"context"
	"fmt"
	"mall/internal/ddd"
	"mall/ordering/internal/domain"
)

type CreateOrderRequest struct {
	ID         string
	CustomerID string
	PaymentID  string
	Items      []domain.Item
}

type CreateOrderHandler struct {
	orders    domain.OrderRepository
	publisher ddd.EventPublisher[ddd.Event]
}

func NewCreateOrderHandler(orders domain.OrderRepository, publisher ddd.EventPublisher[ddd.Event]) CreateOrderHandler {
	return CreateOrderHandler{
		orders:    orders,
		publisher: publisher,
	}
}

func (h CreateOrderHandler) CreateOrder(ctx context.Context, cmd CreateOrderRequest) error {
	order, err := h.orders.Load(ctx, cmd.ID)
	if err != nil {
		return err
	}

	event, err := order.CreateOrder(cmd.ID, cmd.CustomerID, cmd.PaymentID, cmd.Items)
	if err != nil {
		return fmt.Errorf("create order command: %w", err)
	}

	if err := h.orders.Save(ctx, order); err != nil {
		return fmt.Errorf("order creation: %w", err)
	}

	return h.publisher.Publish(ctx, event)
}
