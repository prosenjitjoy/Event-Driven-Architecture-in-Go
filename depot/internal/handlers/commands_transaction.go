package handlers

import (
	"context"
	"database/sql"
	"mall/internal/am"
	"mall/internal/ddd"
	"mall/internal/di"
	"mall/internal/registry"
)

func RegisterCommandHandlersTx(container di.Container) error {
	cmdMsgHandlers := am.RawMessageHandlerFunc(func(ctx context.Context, msg am.IncomingRawMessage) (err error) {
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
		replyStream := di.Get(ctx, "replyStream").(am.ReplyStream)
		commandhandlers := di.Get(ctx, "commandHandlers").(ddd.CommandHandler[ddd.Command])

		cmdHandlers := am.RawMessageHandlerWithMiddleware(
			am.NewCommandMessageHandler(reg, replyStream, commandhandlers),
			di.Get(ctx, "inboxMiddleware").(am.RawMessageHandlerMiddleware),
		)

		return cmdHandlers.HandleMessage(ctx, msg)
	})

	subscriber := container.Get("stream").(am.RawMessageStream)

	return RegisterCommandHandlers(subscriber, cmdMsgHandlers)
}
