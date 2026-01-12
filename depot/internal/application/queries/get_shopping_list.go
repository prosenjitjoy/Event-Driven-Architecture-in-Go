package queries

import (
	"context"
	"mall/depot/internal/domain"
)

type GetShoppingListRequest struct {
	ID string
}

type GetShoppingListHandler struct {
	shoppingList domain.ShoppingListRepository
}

func NewGetShoppingListHandler(shoppingLists domain.ShoppingListRepository) GetShoppingListHandler {
	return GetShoppingListHandler{
		shoppingList: shoppingLists,
	}
}

func (h GetShoppingListHandler) GetShoppingList(ctx context.Context, query GetShoppingListRequest) (*domain.ShoppingList, error) {
	return h.shoppingList.Find(ctx, query.ID)
}
