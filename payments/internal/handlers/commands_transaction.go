package handlers

import (
	"context"
	"database/sql"
	"mall/internal/am"
	"mall/internal/di"
	"mall/payments/internal/constants"
)

func RegisterCommandHandlersTx(container di.Container) error {
	cmdMsgHandlers := am.MessageHandlerFunc(func(ctx context.Context, msg am.IncomingMessage) (err error) {
		ctx = container.Scoped(ctx)

		defer func(tx *sql.Tx) {
			if p := recover(); p != nil {
				tx.Rollback()
				panic(p)
			} else if err != nil {
				tx.Rollback()
			} else {
				err = tx.Commit()
			}
		}(di.Get(ctx, "tx").(*sql.Tx))

		cmdHandlers := di.Get(ctx, constants.CommandHandlersKey).(am.MessageHandler)

		return cmdHandlers.HandleMessage(ctx, msg)
	})

	subscriber := container.Get(constants.MessageSubscriberKey).(am.MessageSubscriber)

	return RegisterCommandHandlers(subscriber, cmdMsgHandlers)
}
