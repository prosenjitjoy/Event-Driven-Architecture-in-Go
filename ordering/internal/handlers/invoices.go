package handlers

import (
	"mall/internal/ddd"
	"mall/ordering/internal/application"
	"mall/ordering/internal/domain"
)

func RegisterInvoiceHandlers(invoiceHandlers application.DomainEventHandlers, domainSubscriber ddd.EventSubscriber) {
	domainSubscriber.Subscribe(domain.OrderReadied{}, invoiceHandlers.OnOrderReadied)
}
