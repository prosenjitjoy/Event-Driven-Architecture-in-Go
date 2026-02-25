package commands

import (
	"context"
	"fmt"
	"mall/depot/internal/domain"
	"mall/internal/ddd"
)

type CreateShoppingListRequest struct {
	ID      string
	OrderID string
	Items   []OrderItem
}

type CreateShoppingListHandler struct {
	shoppingLists   domain.ShoppingListRepository
	stores          domain.StoreRepository
	products        domain.ProductRepository
	domainPublisher ddd.EventPublisher[ddd.AggregateEvent]
}

func NewCreateShoppingListHandler(shoppingLists domain.ShoppingListRepository, stores domain.StoreRepository, products domain.ProductRepository, domainPublisher ddd.EventPublisher[ddd.AggregateEvent]) CreateShoppingListHandler {
	return CreateShoppingListHandler{
		shoppingLists:   shoppingLists,
		stores:          stores,
		products:        products,
		domainPublisher: domainPublisher,
	}
}

func (h CreateShoppingListHandler) CreateShoppingList(ctx context.Context, cmd CreateShoppingListRequest) error {
	list := domain.CreateShopping(cmd.ID, cmd.OrderID)

	for _, item := range cmd.Items {
		store, err := h.stores.Find(ctx, item.StoreID)
		if err != nil {
			return fmt.Errorf("building shopping list: %w", err)
		}
		product, err := h.products.Find(ctx, item.ProductID)
		if err != nil {
			return fmt.Errorf("building shopping list: %w", err)
		}
		err = list.AddItem(store, product, item.Quantity)
		if err != nil {
			return fmt.Errorf("building shopping list: %w", err)
		}
	}

	if err := h.shoppingLists.Save(ctx, list); err != nil {
		return fmt.Errorf("scheduling shopping: %w", err)
	}

	if err := h.domainPublisher.Publish(ctx, list.GetEvents()...); err != nil {
		return err
	}

	return nil
}
