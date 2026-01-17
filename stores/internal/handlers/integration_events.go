package handlers

import (
	"mall/internal/ddd"
	"mall/stores/internal/domain"
)

func RegisterIntegrationEventHandlers(eventHandlers ddd.EventHandler[ddd.AggregateEvent], domainSubscriber ddd.EventSubscriber[ddd.AggregateEvent]) {
	domainSubscriber.Subscribe(domain.StoreCreatedEvent, eventHandlers)
	domainSubscriber.Subscribe(domain.StoreParticipationEnabledEvent, eventHandlers)
	domainSubscriber.Subscribe(domain.StoreParticipationDisabledEvent, eventHandlers)
	domainSubscriber.Subscribe(domain.StoreRebrandedEvent, eventHandlers)
}
