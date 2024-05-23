package application

import (
	"context"
	"mall/internal/ddd"
	"mall/ordering/internal/domain"
)

type ShoppingHandlers struct {
	shopping domain.ShoppingRepository
	ignoreUnimplementedDomainEvents
}

func NewShoppingHandlers(shopping domain.ShoppingRepository) *ShoppingHandlers {
	return &ShoppingHandlers{shopping: shopping}
}

func (h ShoppingHandlers) OnOrderCreated(ctx context.Context, event ddd.Event) error {
	orderCreated := event.(*domain.OrderCreated)

	shoppingID, err := h.shopping.Create(ctx, orderCreated.Order)
	if err != nil {
		return err
	}

	orderCreated.Order.ShoppingID = shoppingID
	return nil
}

func (h ShoppingHandlers) OrderCanceled(ctx context.Context, event ddd.Event) error {
	orderCanceled := event.(*domain.OrderCanceled)

	return h.shopping.Cancel(ctx, orderCanceled.Order.ShoppingID)
}
