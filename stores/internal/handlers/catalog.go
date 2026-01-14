package handlers

import (
	"mall/internal/ddd"
	"mall/stores/internal/domain"
)

func RegisterCatalogHandlers(catalogHandlers ddd.EventHandler[ddd.AggregateEvent], domainSubscriber ddd.EventSubscriber[ddd.AggregateEvent]) {
	domainSubscriber.Subscribe(domain.ProductAddedEvent, catalogHandlers)
	domainSubscriber.Subscribe(domain.ProductRebrandedEvent, catalogHandlers)
	domainSubscriber.Subscribe(domain.ProductPriceIncreaseEvent, catalogHandlers)
	domainSubscriber.Subscribe(domain.ProductPriceDecreaseEvent, catalogHandlers)
	domainSubscriber.Subscribe(domain.ProductRemovedEvent, catalogHandlers)
}
