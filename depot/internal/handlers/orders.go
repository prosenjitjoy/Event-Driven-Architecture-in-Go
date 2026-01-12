package handlers

import (
	"mall/depot/internal/application"
	"mall/depot/internal/domain"
	"mall/internal/ddd"
)

func RegisterOrderHandlers(orderHandler application.DomainEventHandlers, domainSubscriber ddd.EventSubscriber) {
	domainSubscriber.Subscribe(domain.ShoppingListCompleted{}, orderHandler.OnShoppingListCompleted)
}
