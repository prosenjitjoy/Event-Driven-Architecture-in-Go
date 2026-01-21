package handlers

import (
	"mall/internal/ddd"
	"mall/payments/internal/domain"
)

func RegisterIntegrationEventHandlers(eventHandlers ddd.EventHandler[ddd.Event], domainSubscriber ddd.EventSubscriber[ddd.Event]) {
	domainSubscriber.Subscribe(domain.InvoicePaidEvent, eventHandlers)
}
