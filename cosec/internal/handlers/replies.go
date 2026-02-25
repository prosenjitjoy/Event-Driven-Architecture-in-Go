package handlers

import (
	"mall/cosec/internal/domain"
	"mall/internal/am"
	"mall/internal/registry"
	"mall/internal/sec"
)

func NewReplyHandlers(reg registry.Registry, orchestrator sec.Orchestrator[*domain.CreateOrderData], mws ...am.MessageHandlerMiddleware) am.MessageHandler {
	return am.NewReplyHandler(reg, orchestrator, mws...)
}

func RegisterReplyHandlers(subscriber am.MessageSubscriber, handlers am.MessageHandler) error {
	_, err := subscriber.Subscribe(domain.CreateOrderReplyChannel, handlers, am.GroupName("cosec-replies"))
	if err != nil {
		return err
	}

	return nil
}
