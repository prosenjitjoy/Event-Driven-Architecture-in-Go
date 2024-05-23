package handlers

import (
	"mall/internal/ddd"
	"mall/ordering/internal/application"
	"mall/ordering/internal/domain"
)

func RegisterPaymentHandlers(paymentHandlers application.DomainEventHandlers, domainSubscriber ddd.EventSubscriber) {
	domainSubscriber.Subscribe(domain.OrderCreated{}, paymentHandlers.OnOrderCreated)
}
