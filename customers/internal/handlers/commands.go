package handlers

import (
	"context"
	"mall/customers/customerspb"
	"mall/customers/internal/application"
	"mall/internal/am"
	"mall/internal/ddd"
)

type commandHandlers struct {
	app application.App
}

func NewCommandHandlers(app application.App) ddd.CommandHandler[ddd.Command] {
	return commandHandlers{
		app: app,
	}
}

func RegisterCommandHandlers(subscriber am.CommandSubscriber, handlers ddd.CommandHandler[ddd.Command]) error {
	cmdMsgHandler := am.CommandMessageHandlerFunc(func(ctx context.Context, cmdMsg am.IncomingCommandMessage) (ddd.Reply, error) {
		return handlers.HandleCommand(ctx, cmdMsg)
	})

	return subscriber.Subscribe(customerspb.CommandChannel, cmdMsgHandler, am.MessageFilters{
		customerspb.AuthorizeCustomerCommand,
	}, am.GroupName("customer-commands"))
}

func (h commandHandlers) HandleCommand(ctx context.Context, cmd ddd.Command) (ddd.Reply, error) {
	switch cmd.CommandName() {
	case customerspb.AuthorizeCustomerCommand:
		return h.doAuthorizeCustomer(ctx, cmd)
	}

	return nil, nil
}

func (h commandHandlers) doAuthorizeCustomer(ctx context.Context, cmd ddd.Command) (ddd.Reply, error) {
	payload := cmd.Payload().(*customerspb.AuthorizeCustomer)

	if err := h.app.AuthorizeCustomer(ctx, application.AuthorizeCustomer{ID: payload.GetId()}); err != nil {
		return nil, err
	}

	return nil, nil
}
