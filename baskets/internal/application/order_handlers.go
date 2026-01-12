package application

import (
	"context"
	"mall/baskets/internal/domain"
	"mall/internal/ddd"
)

type OrderHandlers struct {
	orders domain.OrderRepository
	ignoreUnimplementedDomainEvents
}

var _ DomainEventHandlers = (*OrderHandlers)(nil)

func NewOrderHandlers(orders domain.OrderRepository) OrderHandlers {
	return OrderHandlers{
		orders: orders,
	}
}

func (h OrderHandlers) OnBasketCheckedOut(ctx context.Context, event ddd.Event) error {
	checkoutOut := event.(*domain.BasketCheckedOut)
	_, err := h.orders.Save(ctx, checkoutOut.Basket)

	return err
}
