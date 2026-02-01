package handlers

import (
	"context"
	"database/sql"
	"mall/baskets/basketspb"
	"mall/depot/depotpb"
	"mall/internal/am"
	"mall/internal/ddd"
	"mall/internal/di"
	"mall/internal/registry"
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

		eventHandler := am.RawMessageHandlerWithMiddleware(
			am.NewEventMessageHandler(reg, integrationEventHandlers),
			di.Get(ctx, "inboxMiddleware").(am.RawMessageHandlerMiddleware),
		)

		return eventHandler.HandleMessage(ctx, msg)
	})

	subscriber := container.Get("stream").(am.RawMessageStream)

	err := subscriber.Subscribe(basketspb.BasketAggregateChannel, eventMsgHandler, am.MessageFilters{
		basketspb.BasketCheckedOutEvent,
	}, am.GroupName("ordering-baskets"))
	if err != nil {
		return err
	}

	err = subscriber.Subscribe(depotpb.ShoppingListAggregateChannel, eventMsgHandler, am.MessageFilters{
		depotpb.ShoppingListCompletedEvent,
	}, am.GroupName("ordering-depot"))
	if err != nil {
		return err
	}

	return nil
}
