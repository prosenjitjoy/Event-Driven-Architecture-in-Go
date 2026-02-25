package domain

import (
	"mall/internal/registry"
	"mall/internal/registry/serdes"
)

const (
	CreateOrderSagaName     = "cosec.CreateOrder"
	CreateOrderReplyChannel = "mall.cosec.replies.CreateOrder"
)

func Registrations(reg registry.Registry) error {
	serde := serdes.NewJsonSerde(reg)

	// saga
	if err := serde.RegisterKey(CreateOrderSagaName, CreateOrderData{}); err != nil {
		return err
	}

	return nil
}
