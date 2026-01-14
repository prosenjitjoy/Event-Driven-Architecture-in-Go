package commands

import (
	"context"
	"mall/depot/internal/domain"
	"mall/internal/ddd"
)

type AssignShoppingListRequest struct {
	ID    string
	BotID string
}

type AssignShoppingListHandler struct {
	shoppingLists   domain.ShoppingListRepository
	domainPublisher ddd.EventPublisher[ddd.AggregateEvent]
}

func NewAssignShoppingListHandler(shoppingList domain.ShoppingListRepository, domainPublisher ddd.EventPublisher[ddd.AggregateEvent]) AssignShoppingListHandler {
	return AssignShoppingListHandler{
		shoppingLists:   shoppingList,
		domainPublisher: domainPublisher,
	}
}

func (h AssignShoppingListHandler) AssignShoppingList(ctx context.Context, cmd AssignShoppingListRequest) error {
	list, err := h.shoppingLists.Find(ctx, cmd.ID)
	if err != nil {
		return err
	}

	if err = list.Assign(cmd.BotID); err != nil {
		return err
	}

	if err = h.shoppingLists.Update(ctx, list); err != nil {
		return err
	}

	if err = h.domainPublisher.Publish(ctx, list.Events()...); err != nil {
		return err
	}

	return nil
}
