package orderingpb

import (
	"mall/internal/registry"
	"mall/internal/registry/serdes"
)

const (
	OrderAggregateChannel = "mall.ordering.events.Order"

	OrderCreatedEvent   = "ordersapi.OrderCreatedEvent"
	OrderRejectedEvent  = "ordersapi.OrderRejectedEvent"
	OrderApprovedEvent  = "ordersapi.OrderApprovedEvent"
	OrderReadiedEvent   = "ordersapi.OrderReadiedEvent"
	OrderCanceledEvent  = "ordersapi.OrderCanceledEvent"
	OrderCompletedEvent = "ordersapi.OrderCompletedEvent"

	CommandChannel = "mall.ordering.commands"

	RejectOrderCommand  = "ordersapi.RejectOrderCommand"
	ApproveOrderCommand = "ordersapi.ApproveOrderCommand"
)

func Registrations(reg registry.Registry) error {
	serde := serdes.NewProtoSerde(reg)

	// Events
	if err := serde.Register(&OrderCreated{}); err != nil {
		return err
	}
	if err := serde.Register(&OrderRejected{}); err != nil {
		return err
	}
	if err := serde.Register(&OrderApproved{}); err != nil {
		return err
	}
	if err := serde.Register(&OrderReadied{}); err != nil {
		return err
	}
	if err := serde.Register(&OrderCanceled{}); err != nil {
		return err
	}
	if err := serde.Register(&OrderCompleted{}); err != nil {
		return err
	}

	// Commands
	if err := serde.Register(&RejectOrder{}); err != nil {
		return err
	}
	if err := serde.Register(&ApproveOrder{}); err != nil {
		return err
	}

	return nil
}

// Events
func (*OrderCreated) Key() string   { return OrderCreatedEvent }
func (*OrderRejected) Key() string  { return OrderRejectedEvent }
func (*OrderApproved) Key() string  { return OrderApprovedEvent }
func (*OrderReadied) Key() string   { return OrderReadiedEvent }
func (*OrderCanceled) Key() string  { return OrderCanceledEvent }
func (*OrderCompleted) Key() string { return OrderCompletedEvent }

// Commands
func (*RejectOrder) Key() string  { return RejectOrderCommand }
func (*ApproveOrder) Key() string { return ApproveOrderCommand }
