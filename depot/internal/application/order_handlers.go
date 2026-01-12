package application

import (
	"context"
	"mall/depot/internal/domain"
	"mall/internal/ddd"
)

type OrderHandlers struct {
	order domain.OrderRepository
	ignoreUnimplementedDomainEvents
}

var _ DomainEventHandlers = (*OrderHandlers)(nil)

func NewOrderHandler(order domain.OrderRepository) OrderHandlers {
	return OrderHandlers{
		order: order,
	}
}

func (h OrderHandlers) OnShoppingListCompleted(ctx context.Context, event ddd.Event) error {
	completed := event.(*domain.ShoppingListCompleted)
	return h.order.Ready(ctx, completed.ShoppingList.OrderID)
}
