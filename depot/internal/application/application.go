package application

import (
	"context"
	"mall/depot/internal/application/commands"
	"mall/depot/internal/application/queries"
	"mall/depot/internal/domain"
	"mall/internal/ddd"
)

type Commands interface {
	CreateShoppingList(ctx context.Context, cmd commands.CreateShoppingListRequest) error
	CancelShoppingList(ctx context.Context, cmd commands.CancelShoppingListRequest) error
	AssignShoppingList(ctx context.Context, cmd commands.AssignShoppingListRequest) error
	CompleteShoppingList(ctx context.Context, cmd commands.CompleteShoppingListRequest) error
}

type Queries interface {
	GetShoppingList(ctx context.Context, query queries.GetShoppingListRequest) (*domain.ShoppingList, error)
}

type App interface {
	Commands
	Queries
}

type appCommands struct {
	commands.CreateShoppingListHandler
	commands.CancelShoppingListHandler
	commands.AssignShoppingListHandler
	commands.CompleteShoppingListHandler
}

type appQueries struct {
	queries.GetShoppingListHandler
}

type Application struct {
	appCommands
	appQueries
}

var _ App = (*Application)(nil)

func New(shoppingLists domain.ShoppingListRepository, stores domain.StoreRepository, products domain.ProductRepository, domainPublisher ddd.EventPublisher[ddd.AggregateEvent]) *Application {
	return &Application{
		appCommands: appCommands{
			CreateShoppingListHandler:   commands.NewCreateShoppingListHandler(shoppingLists, stores, products, domainPublisher),
			CancelShoppingListHandler:   commands.NewCancelShoppingListHandler(shoppingLists, domainPublisher),
			AssignShoppingListHandler:   commands.NewAssignShoppingListHandler(shoppingLists, domainPublisher),
			CompleteShoppingListHandler: commands.NewCompleteShoppingListHandler(shoppingLists, domainPublisher),
		},
		appQueries: appQueries{
			GetShoppingListHandler: queries.NewGetShoppingListHandler(shoppingLists),
		},
	}
}
