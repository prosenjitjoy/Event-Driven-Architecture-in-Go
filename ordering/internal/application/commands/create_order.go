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
	Items      []*domain.Item
}

type CreateOrderHandler struct {
	orders          domain.OrderRepository
	customers       domain.CustomerRepository
	payments        domain.PaymentRepository
	shopping        domain.ShoppingRepository
	domainPublisher ddd.EventPublisher
}

func NewCreateOrderHandler(orders domain.OrderRepository, customers domain.CustomerRepository, payments domain.PaymentRepository, shopping domain.ShoppingRepository, domainPublisher ddd.EventPublisher) CreateOrderHandler {
	return CreateOrderHandler{
		orders:          orders,
		customers:       customers,
		payments:        payments,
		shopping:        shopping,
		domainPublisher: domainPublisher,
	}
}

func (h CreateOrderHandler) CreateOrder(ctx context.Context, cmd CreateOrderRequest) error {
	order, err := domain.CreateOrder(cmd.ID, cmd.CustomerID, cmd.PaymentID, cmd.Items)
	if err != nil {
		return fmt.Errorf("create order command: %w", err)
	}

	if err = h.customers.Authorize(ctx, order.CustomerID); err != nil {
		return fmt.Errorf("order customer authorization: %w", err)
	}

	if err = h.payments.Confirm(ctx, order.PaymentID); err != nil {
		return fmt.Errorf("order payment confirmation: %w", err)
	}

	if order.ShoppingID, err = h.shopping.Create(ctx, order); err != nil {
		return fmt.Errorf("order shopping scheduling: %w", err)
	}

	if err = h.orders.Save(ctx, order); err != nil {
		return fmt.Errorf("order creation: %w", err)
	}

	if err = h.domainPublisher.Publish(ctx, order.GetEvents()...); err != nil {
		return err
	}

	return nil
}
