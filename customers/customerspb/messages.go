package customerspb

import (
	"mall/internal/registry"
	"mall/internal/registry/serdes"
)

const (
	CustomerAggregateChannel = "mall.customers.events.Customer"

	CustomerRegisteredEvent = "customersapi.CustomerRegisteredEvent"
	CustomerSmsChangedEvent = "customersapi.CustomerSmsChangedEvent"
	CustomerEnabledEvent    = "customersapi.CustomerEnabledEvent"
	CustomerDisabledEvent   = "customersapi.CustomerDisabledEvent"

	CommandChannel = "mall.customers.commands"

	AuthorizeCustomerCommand = "customersapi.AuthorizeCustomerCommand"
)

func Registrations(reg registry.Registry) error {
	serde := serdes.NewProtoSerde(reg)

	// Events
	if err := serde.Register(&CustomerRegistered{}); err != nil {
		return err
	}
	if err := serde.Register(&CustomerSmsChanged{}); err != nil {
		return err
	}
	if err := serde.Register(&CustomerEnabled{}); err != nil {
		return err
	}
	if err := serde.Register(&CustomerDisabled{}); err != nil {
		return err
	}

	// Commands
	if err := serde.Register(&AuthorizeCustomer{}); err != nil {
		return err
	}

	return nil
}

// Events
func (*CustomerRegistered) Key() string { return CustomerRegisteredEvent }
func (*CustomerSmsChanged) Key() string { return CustomerSmsChangedEvent }
func (*CustomerEnabled) Key() string    { return CustomerEnabledEvent }
func (*CustomerDisabled) Key() string   { return CustomerDisabledEvent }

// Commands
func (*AuthorizeCustomer) Key() string { return AuthorizeCustomerCommand }
