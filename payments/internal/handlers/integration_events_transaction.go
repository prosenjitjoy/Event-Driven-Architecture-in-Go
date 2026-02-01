package handlers

import (
	"context"
	"database/sql"
	"mall/internal/am"
	"mall/internal/ddd"
	"mall/internal/di"
	"mall/internal/registry"
	"mall/ordering/orderingpb"
)

func RegisterIntegrationEventHandlersTx(container di.Container) error {
	eventMsgHandler := am.RawMessageHandlerFunc(func(ctx context.Context, msg am.IncomingRawMessage) (err error) {
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

		reg := di.Get(ctx, "registry").(registry.Registry)
		integrationEventHandlers := di.Get(ctx, "integrationEventHandlers").(ddd.EventHandler[ddd.Event])

		eventHandlers := am.RawMessageHandlerWithMiddleware(
			am.NewEventMessageHandler(reg, integrationEventHandlers),
			di.Get(ctx, "inboxMiddleware").(am.RawMessageHandlerMiddleware),
		)

		return eventHandlers.HandleMessage(ctx, msg)
	})

	subscriber := container.Get("stream").(am.RawMessageStream)

	err := subscriber.Subscribe(orderingpb.OrderAggregateChannel, eventMsgHandler, am.MessageFilters{
		orderingpb.OrderReadiedEvent,
	}, am.GroupName("payment-orders"))
	if err != nil {
		return err
	}

	return nil
}
