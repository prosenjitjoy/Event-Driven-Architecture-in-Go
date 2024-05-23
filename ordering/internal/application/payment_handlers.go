package application

import (
	"context"
	"mall/internal/ddd"
	"mall/ordering/internal/domain"
)

type PaymentHandlers struct {
	payments domain.PaymentRepository
	ignoreUnimplementedDomainEvents
}

func NewPaymentHandlers(payments domain.PaymentRepository) *PaymentHandlers {
	return &PaymentHandlers{payments: payments}
}

func (h PaymentHandlers) OnOrderCreated(ctx context.Context, event ddd.Event) error {
	orderCreated := event.(*domain.OrderCreated)

	return h.payments.Confirm(ctx, orderCreated.Order.PaymentID)
}
