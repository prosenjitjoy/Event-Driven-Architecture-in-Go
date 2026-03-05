package handlers

import (
	"context"
	"mall/customers/customerspb"
	"mall/customers/internal/application"
	"mall/internal/am"
	"mall/internal/ddd"
	"time"

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
	_, err := subscriber.Subscribe(customerspb.CommandChannel, handlers, am.MessageFilters{
		customerspb.AuthorizeCustomerCommand,
	}, am.GroupName("customer-commands"))
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
