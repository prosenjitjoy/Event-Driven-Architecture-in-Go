package handlers

import (
	"mall/customers/internal/domain"
	"mall/internal/ddd"
)

func RegisterIntegrationEventHandlers(eventHandlers ddd.EventHandler[ddd.AggregateEvent], domainSubscriber ddd.EventSubscriber[ddd.AggregateEvent]) {
	domainSubscriber.Subscribe(domain.CustomerRegisteredEvent, eventHandlers)
	domainSubscriber.Subscribe(domain.CustomerSmsChangedEvent, eventHandlers)
	domainSubscriber.Subscribe(domain.CustomerEnabledEvent, eventHandlers)
	domainSubscriber.Subscribe(domain.CustomerDisabledEvent, eventHandlers)
}
