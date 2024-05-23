package application

import (
	"context"
	"mall/internal/ddd"
	"mall/ordering/internal/domain"
)

type CustomerHandlers struct {
	customers domain.CustomerRepository
	ignoreUnimplementedDomainEvents
}

func NewCustomerRepository(customers domain.CustomerRepository) *CustomerHandlers {
	return &CustomerHandlers{customers: customers}
}

func (h CustomerHandlers) OnOrderCreated(ctx context.Context, event ddd.Event) error {
	orderCreated := event.(*domain.OrderCreated)

	return h.customers.Authorize(ctx, orderCreated.Order.CustomerID)
}
