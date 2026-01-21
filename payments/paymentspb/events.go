package paymentspb

import (
	"mall/internal/registry"
	"mall/internal/registry/serdes"
)

const (
	InvoiceAggregateChannel = "mall.payments.events.Invoice"

	InvoicePaidEvent = "paymentsapi.InvoicePaid"
)

func Registrations(reg registry.Registry) error {
	serde := serdes.NewProtoSerde(reg)

	// invoice events
	if err := serde.Register(&InvoicePaid{}); err != nil {
		return err
	}

	return nil
}

func (*InvoicePaid) Key() string { return InvoicePaidEvent }
