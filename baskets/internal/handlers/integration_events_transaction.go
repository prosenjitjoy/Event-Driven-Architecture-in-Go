package handlers

import (
	"context"
	"database/sql"
	"mall/internal/am"
	"mall/internal/ddd"
	"mall/internal/di"
	"mall/internal/registry"
	"mall/stores/storespb"
)

func RegisterIntegrationEventHandlersTx(container di.Container) error {
	eventMsgHandler := am.RawMessageHandlerFunc(func(ctx context.Context, msg am.IncomingRawMessage) (err error) {
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
		}(di.Get(ctx, "tx").(*sql.Tx))

		reg := di.Get(ctx, "registry").(registry.Registry)
		integrationEventHandlers := di.Get(ctx, "integrationEventHandlers").(ddd.EventHandler[ddd.Event])

		evtMsgHandlers := am.RawMessageHandlerWithMiddleware(
			am.NewEventMessageHandler(reg, integrationEventHandlers),
			di.Get(ctx, "inboxMiddleware").(am.RawMessageHandlerMiddleware),
		)

		return evtMsgHandlers.HandleMessage(ctx, msg)
	})

	subscriber := container.Get("stream").(am.RawMessageStream)

	err := subscriber.Subscribe(storespb.StoreAggregateChannel, eventMsgHandler, am.MessageFilters{
		storespb.StoreCreatedEvent,
		storespb.StoreRebrandedEvent,
	}, am.GroupName("baskets-stores"))
	if err != nil {
		return err
	}

	err = subscriber.Subscribe(storespb.ProductAggregateChannel, eventMsgHandler, am.MessageFilters{
		storespb.ProductAddedEvent,
		storespb.ProductRebrandedEvent,
		storespb.ProductPriceIncreasedEvent,
		storespb.ProductPriceDecreasedEvent,
		storespb.ProductRemovedEvent,
	}, am.GroupName("baskets-products"))
	if err != nil {
		return err
	}

	return nil
}
