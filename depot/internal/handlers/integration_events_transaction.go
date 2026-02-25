package handlers

import (
	"context"
	"database/sql"
	"mall/customers/constants"
	"mall/internal/am"
	"mall/internal/di"
)

func RegisterIntegrationEventHandlersTx(container di.Container) error {
	eventMsgHandler := am.MessageHandlerFunc(func(ctx context.Context, msg am.IncomingMessage) (err error) {
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

		eventHandler := di.Get(ctx, constants.IntegrationEventHandlersKey).(am.MessageHandler)

		return eventHandler.HandleMessage(ctx, msg)
	})

	subscriber := container.Get(constants.MessageSubscriberKey).(am.MessageSubscriber)

	return RegisterIntegrationEventHandlers(subscriber, eventMsgHandler)
}
