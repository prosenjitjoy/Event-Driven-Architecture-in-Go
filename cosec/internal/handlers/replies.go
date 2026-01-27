package handlers

import (
	"context"
	"mall/cosec/internal/domain"
	"mall/internal/am"
	"mall/internal/sec"
)

func RegisterReplyHandlers(subscriber am.ReplySubscriber, orchestrator sec.Orchestrator[*domain.CreateOrderData]) error {
	replyMsgHandler := am.MessageHandlerFunc[am.IncomingReplyMessage](func(ctx context.Context, replyMsg am.IncomingReplyMessage) error {
		return orchestrator.HandleReply(ctx, replyMsg)
	})

	return subscriber.Subscribe(orchestrator.ReplyTopic(), replyMsgHandler, am.GroupName("cosec-replies"))
}
