package handlers

import (
	"context"
	"database/sql"
	"mall/cosec/internal/constants"
	"mall/internal/am"
	"mall/internal/di"
)

func RegisterReplyHandlersTx(container di.Container) error {
	replyMsgHandler := am.MessageHandlerFunc(func(ctx context.Context, msg am.IncomingMessage) (err error) {
		ctx = container.Scoped(ctx)

		defer func(tx *sql.Tx) {
			if p := recover(); p != nil {
				tx.Rollback()
				panic(p)
			} else if err != nil {

			}
		}(di.Get(ctx, constants.DatabaseTransactionKey).(*sql.Tx))

		replyHandlers := di.Get(ctx, constants.ReplyHandlersKey).(am.MessageHandler)

		return replyHandlers.HandleMessage(ctx, msg)
	})

	subscriber := container.Get(constants.MessageSubscriberKey).(am.MessageSubscriber)

	return RegisterReplyHandlers(subscriber, replyMsgHandler)
}
