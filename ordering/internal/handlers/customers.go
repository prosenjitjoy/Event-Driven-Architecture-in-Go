package handlers

import (
	"mall/internal/ddd"
	"mall/ordering/internal/application"
	"mall/ordering/internal/domain"
)

func RegisterCustomerHandlers(customerHandlers application.DomainEventHandlers, domainSubscriber ddd.EventSubscriber) {
	domainSubscriber.Subscribe(domain.OrderCreated{}, customerHandlers.OnOrderCreated)
}
