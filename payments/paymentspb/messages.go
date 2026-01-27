package paymentspb

import (
	"mall/internal/registry"
	"mall/internal/registry/serdes"
)

const (
	InvoiceAggregateChannel = "mall.payments.events.Invoice"

	InvoicePaidEvent = "paymentsapi.InvoicePaidEvent"

	CommandChannel = "mall.payments.commands"

	ConfirmPaymentCommand = "paymentsapi.ConfirmPaymentCommand"
)

func Registrations(reg registry.Registry) error {
	serde := serdes.NewProtoSerde(reg)

	// Events
	if err := serde.Register(&InvoicePaid{}); err != nil {
		return err
	}

	// Commands
	if err := serde.Register(&ConfirmPayment{}); err != nil {
		return err
	}

	return nil
}

// Events
func (*InvoicePaid) Key() string { return InvoicePaidEvent }

// Commands
func (*ConfirmPayment) Key() string { return ConfirmPaymentCommand }
