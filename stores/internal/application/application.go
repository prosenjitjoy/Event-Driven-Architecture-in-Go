package application

import (
	"context"
	"mall/stores/internal/application/commands"
	"mall/stores/internal/application/queries"
	"mall/stores/internal/domain"
)

type Commands interface {
	CreateStore(ctx context.Context, cmd commands.CreateStoreRequest) error
	EnableParticipation(ctx context.Context, cmd commands.EnableParticipationRequest) error
	DisableParticipation(ctx context.Context, cmd commands.DisableParticipationRequest) error
	RebrandStore(ctx context.Context, cmd commands.RebrandStoreRequest) error
	AddProduct(ctx context.Context, cmd commands.AddProductRequest) error
	RebrandProduct(ctx context.Context, cmd commands.RebrandProductRequest) error
	IncreaseProductPrice(ctx context.Context, cmd commands.IncreaseProductPriceRequest) error
	DecreaseProductPrice(ctx context.Context, cmd commands.DecreaseProductPriceRequest) error
	RemoveProduct(ctx context.Context, cmd commands.RemoveProductRequest) error
}

type Queries interface {
	GetStore(ctx context.Context, query queries.GetStoreRequest) (*domain.MallStore, error)
	GetStores(ctx context.Context, query queries.GetStoresRequest) ([]*domain.MallStore, error)
	GetParticipatingStores(ctx context.Context, query queries.GetParticipatingStoreRequest) ([]*domain.MallStore, error)
	GetCatalog(ctx context.Context, query queries.GetCatalogRequest) ([]*domain.CatalogProduct, error)
	GetProduct(ctx context.Context, query queries.GetProductRequest) (*domain.CatalogProduct, error)
}

type App interface {
	Commands
	Queries
}

type appCommands struct {
	commands.CreateStoreHandler
	commands.EnableParticipationHandler
	commands.DisableParticipationHandler
	commands.RebrandStoreHandler
	commands.AddProductHandler
	commands.RebrandProductHandler
	commands.IncreaseProductPriceHandler
	commands.DecreaseProductPriceHandler
	commands.RemoveProductHandler
}

type appQueries struct {
	queries.GetStoreHandler
	queries.GetStoresHandler
	queries.GetParticipatingStoreHandler
	queries.GetCatalogHandler
	queries.GetProductHandler
}

type Application struct {
	appCommands
	appQueries
}

var _ App = (*Application)(nil)

func New(stores domain.StoreRepository, products domain.ProductRepository, catalog domain.CatalogRepository, mall domain.MallRepository) *Application {
	return &Application{
		appCommands: appCommands{
			CreateStoreHandler:          commands.NewCreateStoreHandler(stores),
			EnableParticipationHandler:  commands.NewEnableParticipationHandler(stores),
			DisableParticipationHandler: commands.NewDisableParticipationHandler(stores),
			RebrandStoreHandler:         commands.NewRebrandStoreHandler(stores),
			AddProductHandler:           commands.NewAddProductHandler(products),
			RebrandProductHandler:       commands.NewRebrandProductHandler(products),
			IncreaseProductPriceHandler: commands.NewIncreaseProductPriceHandler(products),
			DecreaseProductPriceHandler: commands.NewDecreaseProductPriceHandler(products),
			RemoveProductHandler:        commands.NewRemoveProductHandler(products),
		},
		appQueries: appQueries{
			GetStoreHandler:              queries.NewGetStoreHandler(mall),
			GetStoresHandler:             queries.NewGetStoresHandler(mall),
			GetParticipatingStoreHandler: queries.NewGetParticipatingStoreHandler(mall),
			GetCatalogHandler:            queries.NewGetCatalogHandler(catalog),
			GetProductHandler:            queries.NewGetProductHandler(catalog),
		},
	}
}
