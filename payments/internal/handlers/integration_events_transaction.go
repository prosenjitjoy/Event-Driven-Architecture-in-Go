package handlers

import (
	"context"
	"database/sql"
	"mall/internal/am"
	"mall/internal/di"
	"mall/payments/internal/constants"
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

		eventHandlers := di.Get(ctx, constants.IntegrationEventHandlersKey).(am.MessageHandler)

		return eventHandlers.HandleMessage(ctx, msg)
	})

	subscriber := container.Get(constants.MessageSubscriberKey).(am.MessageSubscriber)

	return RegisterIntegrationEventHandlers(subscriber, eventMsgHandler)
}
