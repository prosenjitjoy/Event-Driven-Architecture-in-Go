package handlers

import (
	"context"
	"mall/internal/am"
	"mall/internal/ddd"
	"mall/payments/internal/application"
	"mall/payments/paymentspb"
)

type commandHandlers struct {
	app application.App
}

func NewCommandHandlers(app application.App) ddd.CommandHandler[ddd.Command] {
	return commandHandlers{app: app}
}

func RegisterCommandHandlers(subscriber am.MessageSubscriber, handlers am.MessageHandler) error {
	_, err := subscriber.Subscribe(paymentspb.CommandChannel, handlers, am.MessageFilters{
		paymentspb.ConfirmPaymentCommand,
	}, am.GroupName("payment-commands"))
	if err != nil {
		return err
	}

	return nil
}

func (h commandHandlers) HandleCommand(ctx context.Context, cmd ddd.Command) (ddd.Reply, error) {
	switch cmd.CommandName() {
	case paymentspb.ConfirmPaymentCommand:
		return h.doConfirmPayment(ctx, cmd)
	}

	return nil, nil
}

func (h commandHandlers) doConfirmPayment(ctx context.Context, cmd ddd.Command) (ddd.Reply, error) {
	payload := cmd.Payload().(*paymentspb.ConfirmPayment)

	err := h.app.ConfirmPayment(ctx, application.ConfirmPayment{ID: payload.GetId()})
	if err != nil {
		return nil, err
	}

	return nil, nil
}
