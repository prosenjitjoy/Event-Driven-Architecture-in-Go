package handlers

import (
	"mall/internal/ddd"
	"mall/stores/internal/domain"
)

func RegisterCatalogHandlers(catalogHandlers ddd.EventHandler[ddd.AggregateEvent], domainSubscriber ddd.EventSubscriber[ddd.AggregateEvent]) {
	domainSubscriber.Subscribe(catalogHandlers,
		domain.ProductAddedEvent,
		domain.ProductRebrandedEvent,
		domain.ProductPriceIncreaseEvent,
		domain.ProductPriceDecreaseEvent,
		domain.ProductRemovedEvent,
	)
}
