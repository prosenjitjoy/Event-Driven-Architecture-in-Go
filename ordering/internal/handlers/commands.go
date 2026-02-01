package handlers

import (
	"context"
	"mall/internal/am"
	"mall/internal/ddd"
	"mall/ordering/internal/application"
	"mall/ordering/internal/application/commands"
	"mall/ordering/orderingpb"
)

type commandHandlers struct {
	app application.App
}

func NewCommandHandlers(app application.App) ddd.CommandHandler[ddd.Command] {
	return commandHandlers{app: app}
}

func RegisterCommandHandlers(subscriber am.RawMessageSubscriber, handlers am.RawMessageHandler) error {
	err := subscriber.Subscribe(orderingpb.CommandChannel, handlers, am.MessageFilters{
		orderingpb.RejectOrderCommand,
		orderingpb.ApproveOrderCommand,
	}, am.GroupName("ordering-commands"))
	if err != nil {
		return err
	}

	return nil
}

func (h commandHandlers) HandleCommand(ctx context.Context, cmd ddd.Command) (ddd.Reply, error) {
	switch cmd.CommandName() {
	case orderingpb.RejectOrderCommand:
		return h.doRejectOrder(ctx, cmd)
	case orderingpb.ApproveOrderCommand:
		return h.doApproveOrder(ctx, cmd)
	}

	return nil, nil
}

func (h commandHandlers) doRejectOrder(ctx context.Context, cmd ddd.Command) (ddd.Reply, error) {
	payload := cmd.Payload().(*orderingpb.RejectOrder)

	err := h.app.RejectOrder(ctx, commands.RejectOrderRequest{
		ID: payload.GetId(),
	})
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func (h commandHandlers) doApproveOrder(ctx context.Context, cmd ddd.Command) (ddd.Reply, error) {
	payload := cmd.Payload().(*orderingpb.ApproveOrder)

	err := h.app.ApproveOrder(ctx, commands.ApproveOrderRequest{
		ID:         payload.GetId(),
		ShoppingID: payload.GetShoppingId(),
	})
	if err != nil {
		return nil, err
	}

	return nil, nil
}
