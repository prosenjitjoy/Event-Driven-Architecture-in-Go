package storespb

import (
	"mall/internal/registry"
	"mall/internal/registry/serdes"
)

const (
	StoreAggregateChannel = "mall.stores.events.Store"

	StoreCreatedEvent              = "storesapi.StoreCreatedEvent"
	StoreParticipationToggledEvent = "storesapi.StoreParticipationToggledEvent"
	StoreRebrandedEvent            = "storesapi.StoreRebrandedEvent"

	ProductAggregateChannel = "mall.stores.events.Product"

	ProductAddedEvent          = "storesapi.ProductAddedEvent"
	ProductRebrandedEvent      = "storesapi.ProductRebrandedEvent"
	ProductPriceIncreasedEvent = "storesapi.ProductPriceIncreasedEvent"
	ProductPriceDecreasedEvent = "storesapi.ProductPriceDecreasedEvent"
	ProductRemovedEvent        = "storesapi.ProductRemovedEvent"
)

func Registrations(reg registry.Registry) error {
	serde := serdes.NewProtoSerde(reg)

	// Events
	if err := serde.Register(&StoreCreated{}); err != nil {
		return err
	}
	if err := serde.Register(&StoreParticipationToggled{}); err != nil {
		return err
	}
	if err := serde.Register(&StoreRebranded{}); err != nil {
		return err
	}

	if err := serde.Register(&ProductAdded{}); err != nil {
		return err
	}
	if err := serde.Register(&ProductRebranded{}); err != nil {
		return err
	}
	if err := serde.RegisterKey(ProductPriceIncreasedEvent, &ProductPriceChanged{}); err != nil {
		return err
	}
	if err := serde.RegisterKey(ProductPriceDecreasedEvent, &ProductPriceChanged{}); err != nil {
		return err
	}
	if err := serde.Register(&ProductRemoved{}); err != nil {
		return err
	}

	return nil
}

// Events
func (*StoreCreated) Key() string              { return StoreCreatedEvent }
func (*StoreParticipationToggled) Key() string { return StoreParticipationToggledEvent }
func (*StoreRebranded) Key() string            { return StoreRebrandedEvent }

func (*ProductAdded) Key() string     { return ProductAddedEvent }
func (*ProductRebranded) Key() string { return ProductRebrandedEvent }
func (*ProductRemoved) Key() string   { return ProductRemovedEvent }
