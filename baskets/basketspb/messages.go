package basketspb

import (
	"mall/internal/registry"
	"mall/internal/registry/serdes"
)

const (
	BasketAggregateChannel = "mall.baskets.events.Basket"

	BasketStartedEvent    = "basketsapi.BasketStartedEvent"
	BasketCanceledEvent   = "basketsapi.BasketCanceledEvent"
	BasketCheckedOutEvent = "basketsapi.BasketCheckedOutEvent"
)

func Registrations(reg registry.Registry) error {
	serde := serdes.NewProtoSerde(reg)

	// Events
	if err := serde.Register(&BasketStarted{}); err != nil {
		return err
	}
	if err := serde.Register(&BasketCanceled{}); err != nil {
		return err
	}
	if err := serde.Register(&BasketCheckedOut{}); err != nil {
		return err
	}

	return nil
}

// Events
func (*BasketStarted) Key() string    { return BasketStartedEvent }
func (*BasketCanceled) Key() string   { return BasketCanceledEvent }
func (*BasketCheckedOut) Key() string { return BasketCheckedOutEvent }
