package application

import (
	"context"
	"mall/internal/ddd"
	"mall/stores/internal/application/commands"
	"mall/stores/internal/application/queries"
	"mall/stores/internal/domain"
)

type Commands interface {
	CreateStore(ctx context.Context, cmd commands.CreateStoreRequest) error
	EnableParticipation(ctx context.Context, cmd commands.EnableParticipationRequest) error
	DisableParticipation(ctx context.Context, cmd commands.DisableParticipationRequest) error
	AddProduct(ctx context.Context, cmd commands.AddProductRequest) error
	RemoveProduct(ctx context.Context, cmd commands.RemoveProductRequest) error
}

type Queries interface {
	GetStore(ctx context.Context, query queries.GetStoreRequest) (*domain.Store, error)
	GetStores(ctx context.Context, query queries.GetStoresRequest) ([]*domain.Store, error)
	GetParticipatingStores(ctx context.Context, query queries.GetParticipatingStoreRequest) ([]*domain.Store, error)
	GetCatalog(ctx context.Context, query queries.GetCatalogRequest) ([]*domain.Product, error)
	GetProduct(ctx context.Context, query queries.GetProductRequest) (*domain.Product, error)
}

type App interface {
	Commands
	Queries
}

type appCommands struct {
	commands.CreateStoreHandler
	commands.EnableParticipationHandler
	commands.DisableParticipationHandler
	commands.AddProductHandler
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

func New(stores domain.StoreRepository, participatingStores domain.ParticipatingStoreRepository, products domain.ProductRepository, domainPublisher ddd.EventPublisher) *Application {
	return &Application{
		appCommands: appCommands{
			CreateStoreHandler:          commands.NewCreateStoreHandler(stores, domainPublisher),
			EnableParticipationHandler:  commands.NewEnableParticipationHandler(stores, domainPublisher),
			DisableParticipationHandler: commands.NewDisableParticipationHandler(stores, domainPublisher),
			AddProductHandler:           commands.NewAddProductHandler(stores, products, domainPublisher),
		},
		appQueries: appQueries{
			GetStoreHandler:              queries.NewGetStoreHandler(stores),
			GetStoresHandler:             queries.NewGetStoresHandler(stores),
			GetParticipatingStoreHandler: queries.NewGetParticipatingStoreHandler(participatingStores),
			GetCatalogHandler:            queries.NewGetCatalogHandler(products),
			GetProductHandler:            queries.NewGetProductHandler(products),
		},
	}
}
