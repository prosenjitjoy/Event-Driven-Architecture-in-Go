package handlers

import (
	"context"
	"mall/depot/depotpb"
	"mall/depot/internal/application"
	"mall/depot/internal/application/commands"
	"mall/internal/am"
	"mall/internal/ddd"
	"time"

	"github.com/google/uuid"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

type commandHandlers struct {
	app application.App
}

var _ ddd.CommandHandler[ddd.Command] = (*commandHandlers)(nil)

func NewCommandHandlers(app application.App) ddd.CommandHandler[ddd.Command] {
	return commandHandlers{app: app}
}

func RegisterCommandHandlers(subscriber am.MessageSubscriber, handlers am.MessageHandler) error {
	_, err := subscriber.Subscribe(depotpb.CommandChannel, handlers, am.MessageFilters{
		depotpb.CreateShoppingListCommand,
		depotpb.CancelShoppingListCommand,
		depotpb.InitiateShoppingCommand,
	}, am.GroupName("depot-commands"))
	if err != nil {
		return err
	}

	return nil
}

func (h commandHandlers) HandleCommand(ctx context.Context, cmd ddd.Command) (reply ddd.Reply, err error) {
	span := trace.SpanFromContext(ctx)

	defer func(started time.Time) {
		if err != nil {
			span.AddEvent(
				"Encountered an error handling command",
				trace.WithAttributes(
					attribute.String("error.message", err.Error()),
				),
			)
			span.SetStatus(codes.Error, err.Error())
		}

		span.AddEvent(
			"Handled command",
			trace.WithAttributes(
				attribute.Int64("TookMS", int64(time.Since(started).Milliseconds())),
			),
		)
	}(time.Now())

	span.AddEvent(
		"Handling command",
		trace.WithAttributes(
			attribute.String("Command", cmd.CommandName()),
		),
	)

	switch cmd.CommandName() {
	case depotpb.CreateShoppingListCommand:
		return h.doCreateShoppingList(ctx, cmd)
	case depotpb.CancelShoppingListCommand:
		return h.doCancelShoppingList(ctx, cmd)
	}

	return nil, nil
}

func (h commandHandlers) doCreateShoppingList(ctx context.Context, cmd ddd.Command) (ddd.Reply, error) {
	payload := cmd.Payload().(*depotpb.CreateShoppingList)

	id := uuid.New().String()

	items := make([]commands.OrderItem, 0, len(payload.GetItems()))

	for _, item := range payload.GetItems() {
		items = append(items, commands.OrderItem{
			StoreID:   item.GetStoreId(),
			ProductID: item.GetProductId(),
			Quantity:  int(item.GetQuantity()),
		})
	}

	err := h.app.CreateShoppingList(ctx, commands.CreateShoppingListRequest{
		ID:      id,
		OrderID: payload.GetOrderId(),
		Items:   items,
	})
	if err != nil {
		return nil, err
	}

	return ddd.NewReply(depotpb.CreatedShoppingListReply, &depotpb.CreatedShoppingList{Id: id}), nil
}

func (h commandHandlers) doCancelShoppingList(ctx context.Context, cmd ddd.Command) (ddd.Reply, error) {
	payload := cmd.Payload().(*depotpb.CancelShoppingList)

	err := h.app.CancelShoppingList(ctx, commands.CancelShoppingListRequest{ID: payload.GetId()})
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func (h commandHandlers) doInitiateShoppingList(ctx context.Context, cmd ddd.Command) (ddd.Reply, error) {
	payload := cmd.Payload().(*depotpb.InitiateShopping)

	err := h.app.InitiateShoppingList(ctx, commands.InitiateShoppingRequest{ID: payload.GetId()})
	if err != nil {
		return nil, err
	}

	return nil, nil
}
