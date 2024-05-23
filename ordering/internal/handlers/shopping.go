package handlers

import (
	"mall/internal/ddd"
	"mall/ordering/internal/application"
	"mall/ordering/internal/domain"
)

func RegisterShoppingHandlers(shoppingHandlers application.DomainEventHandlers, domainSubscriber ddd.EventSubscriber) {
	domainSubscriber.Subscribe(domain.OrderCreated{}, shoppingHandlers.OnOrderCreated)
	domainSubscriber.Subscribe(domain.OrderCanceled{}, shoppingHandlers.OnOrderCanceled)
}
