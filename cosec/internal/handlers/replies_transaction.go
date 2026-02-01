package handlers

import (
	"context"
	"database/sql"
	"mall/cosec/internal/application"
	"mall/cosec/internal/domain"
	"mall/internal/am"
	"mall/internal/di"
	"mall/internal/registry"
	"mall/internal/sec"
)

func RegisterReplyHandlersTx(container di.Container) error {
	replyMsgHandler := am.RawMessageHandlerFunc(func(ctx context.Context, msg am.IncomingRawMessage) (err error) {
		ctx = container.Scoped(ctx)

		defer func(tx *sql.Tx) {
			if p := recover(); p != nil {
				tx.Rollback()
				panic(p)
			} else if err != nil {

			}
		}(di.Get(ctx, "tx").(*sql.Tx))

		reg := di.Get(ctx, "registry").(registry.Registry)
		orchestrator := di.Get(ctx, "orchestrator").(sec.Orchestrator[*domain.CreateOrderData])

		replyHandlers := am.RawMessageHandlerWithMiddleware(
			am.NewReplyMessageHandler(reg, orchestrator),
			di.Get(ctx, "inboxMiddleware").(am.RawMessageHandlerMiddleware),
		)

		return replyHandlers.HandleMessage(ctx, msg)
	})

	subscriber := container.Get("stream").(am.RawMessageStream)

	return subscriber.Subscribe(application.CreateOrderReplyChannel, replyMsgHandler, am.GroupName("cosec-replies"))
}
