package orderingpb

import (
	"mall/internal/registry"
	"mall/internal/registry/serdes"
)

const (
	OrderAggregateChannel = "mall.ordering.events.Order"

	OrderCreatedEvent   = "ordersapi.OrderCreated"
	OrderReadiedEvent   = "ordersapi.OrderReadied"
	OrderCanceledEvent  = "ordersapi.OrderCanceled"
	OrderCompletedEvent = "ordersapi.OrderCompleted"
)

func Registrations(reg registry.Registry) error {
	serde := serdes.NewProtoSerde(reg)

	// order events
	if err := serde.Register(&OrderCreated{}); err != nil {
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

	return nil
}

func (*OrderCreated) Key() string   { return OrderCreatedEvent }
func (*OrderReadied) Key() string   { return OrderReadiedEvent }
func (*OrderCanceled) Key() string  { return OrderCanceledEvent }
func (*OrderCompleted) Key() string { return OrderCompletedEvent }
