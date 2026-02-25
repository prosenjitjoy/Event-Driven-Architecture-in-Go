package domain

import (
	"mall/internal/es"
	"mall/internal/registry"
	"mall/internal/registry/serdes"
)

const (
	OrderCreatedEvent   = "ordering.OrderCreated"
	OrderRejectedEvent  = "ordering.OrderRejected"
	OrderApprovedEvent  = "ordering.OrderApproved"
	OrderCanceledEvent  = "ordering.OrderCanceled"
	OrderReadiedEvent   = "ordering.OrderReadied"
	OrderCompletedEvent = "ordering.OrderCompleted"
)

func Registrations(reg registry.Registry) error {
	serde := serdes.NewJsonSerde(reg)

	// order
	if err := serde.Register(Order{}, func(v any) error {
		order := v.(*Order)
		order.Aggregate = es.NewAggregate("", OrderAggregate)
		return nil
	}); err != nil {
		return err
	}

	// order events
	if err := serde.Register(OrderCreated{}); err != nil {
		return err
	}
	if err := serde.Register(OrderRejected{}); err != nil {
		return err
	}
	if err := serde.Register(OrderApproved{}); err != nil {
		return err
	}
	if err := serde.Register(OrderCanceled{}); err != nil {
		return err
	}
	if err := serde.Register(OrderReadied{}); err != nil {
		return err
	}
	if err := serde.Register(OrderCompleted{}); err != nil {
		return err
	}

	// order snapshots
	if err := serde.RegisterKey(OrderV1{}.SnapshotName(), OrderV1{}); err != nil {
		return err
	}

	return nil
}

func (Order) Key() string { return OrderAggregate }

func (OrderCreated) Key() string   { return OrderCreatedEvent }
func (OrderRejected) Key() string  { return OrderRejectedEvent }
func (OrderApproved) Key() string  { return OrderApprovedEvent }
func (OrderCanceled) Key() string  { return OrderCanceledEvent }
func (OrderReadied) Key() string   { return OrderReadiedEvent }
func (OrderCompleted) Key() string { return OrderCompletedEvent }
