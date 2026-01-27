package commands

import (
	"context"
	"mall/depot/internal/domain"
	"mall/internal/ddd"
)

type InitiateShoppingRequest struct {
	ID string
}

type InitiateShoppingListHandler struct {
	shoppingLists   domain.ShoppingListRepository
	domainPublisher ddd.EventPublisher[ddd.AggregateEvent]
}

func NewInitiateShoppingListHandler(lists domain.ShoppingListRepository, publisher ddd.EventPublisher[ddd.AggregateEvent]) InitiateShoppingListHandler {
	return InitiateShoppingListHandler{
		shoppingLists:   lists,
		domainPublisher: publisher,
	}
}

func (h InitiateShoppingListHandler) InitiateShoppingList(ctx context.Context, cmd InitiateShoppingRequest) error {
	list, err := h.shoppingLists.Find(ctx, cmd.ID)
	if err != nil {
		return err
	}

	if err := list.Initiate(); err != nil {
		return err
	}

	if err := h.shoppingLists.Update(ctx, list); err != nil {
		return err
	}

	// publish domain events
	if err := h.domainPublisher.Publish(ctx, list.Events()...); err != nil {
		return err
	}

	return nil
}
