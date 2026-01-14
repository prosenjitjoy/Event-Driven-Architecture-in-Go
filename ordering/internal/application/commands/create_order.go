package commands

import (
	"context"
	"fmt"
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
	customers domain.CustomerRepository
	payments  domain.PaymentRepository
	shopping  domain.ShoppingRepository
}

func NewCreateOrderHandler(orders domain.OrderRepository, customers domain.CustomerRepository, payments domain.PaymentRepository, shopping domain.ShoppingRepository) CreateOrderHandler {
	return CreateOrderHandler{
		orders:    orders,
		customers: customers,
		payments:  payments,
		shopping:  shopping,
	}
}

func (h CreateOrderHandler) CreateOrder(ctx context.Context, cmd CreateOrderRequest) error {
	order, err := h.orders.Load(ctx, cmd.ID)
	if err != nil {
		return err
	}

	// authorizeCustomer
	if err = h.customers.Authorize(ctx, cmd.CustomerID); err != nil {
		return fmt.Errorf("order customer authorization: %w", err)
	}

	// validatePayment
	if err = h.payments.Confirm(ctx, cmd.PaymentID); err != nil {
		return fmt.Errorf("order payment confirmation: %w", err)
	}

	// scheduleShopping
	shoppingID, err := h.shopping.Create(ctx, cmd.ID, cmd.Items)
	if err != nil {
		return fmt.Errorf("order shopping scheduling: %w", err)
	}

	// order creationg
	err = order.CreateOrder(cmd.ID, cmd.CustomerID, cmd.PaymentID, shoppingID, cmd.Items)
	if err != nil {
		return fmt.Errorf("create order command: %w", err)
	}

	if err = h.orders.Save(ctx, order); err != nil {
		return fmt.Errorf("order creation: %w", err)
	}

	return nil
}
