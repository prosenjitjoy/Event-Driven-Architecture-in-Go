package handlers

import (
	"context"
	"database/sql"
	"mall/baskets/internal/constants"
	"mall/internal/am"
	"mall/internal/di"
)

func RegisterIntegrationEventHandlersTx(container di.Container) error {
	eventMsgHandler := am.MessageHandlerFunc(func(ctx context.Context, msg am.IncomingMessage) (err error) {
		ctx = container.Scoped(ctx)

		defer func(tx *sql.Tx) {
			if p := recover(); p != nil {
				_ = tx.Rollback()
				panic(p)
			} else if err != nil {
				_ = tx.Rollback()
			} else {
				err = tx.Commit()
			}
		}(di.Get(ctx, constants.DatabaseTransactionKey).(*sql.Tx))

		evtHandlers := di.Get(ctx, constants.IntegrationEventHandlersKey).(am.MessageHandler)

		return evtHandlers.HandleMessage(ctx, msg)
	})

	subscriber := container.Get(constants.MessageSubscriberKey).(am.MessageSubscriber)

	return RegisterIntegrationEventHandlers(subscriber, eventMsgHandler)
}
