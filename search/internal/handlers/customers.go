package handlers

import (
	"context"
	"mall/customers/customerspb"
	"mall/internal/am"
	"mall/internal/ddd"
)

func RegisterCustomerHandlers(customerHandlers ddd.EventHandler[ddd.Event], stream am.EventSubscriber) error {
	eventMsgHandler := am.MessageHandlerFunc[am.EventMessage](func(ctx context.Context, eventMsg am.EventMessage) error {
		return customerHandlers.HandleEvent(ctx, eventMsg)
	})

	return stream.Subscribe(customerspb.CustomerAggregateChannel, eventMsgHandler, am.MessageFilters{
		customerspb.CustomerRegisteredEvent,
	}, am.GroupName("search-customers"))
}
