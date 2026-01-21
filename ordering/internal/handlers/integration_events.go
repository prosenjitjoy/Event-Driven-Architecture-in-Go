package handlers

import (
	"mall/internal/ddd"
	"mall/ordering/internal/domain"
)

func RegisterIntegrationEventHandlers(eventHandlers ddd.EventHandler[ddd.AggregateEvent], domainSubscriber ddd.EventSubscriber[ddd.AggregateEvent]) {
	domainSubscriber.Subscribe(domain.OrderCreatedEvent, eventHandlers)
	domainSubscriber.Subscribe(domain.OrderReadiedEvent, eventHandlers)
	domainSubscriber.Subscribe(domain.OrderCanceledEvent, eventHandlers)
	domainSubscriber.Subscribe(domain.OrderCompletedEvent, eventHandlers)
}
