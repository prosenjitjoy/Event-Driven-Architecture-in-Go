package logging

import (
	"context"
	"log/slog"
	"mall/stores/internal/application"
	"mall/stores/internal/application/commands"
	"mall/stores/internal/application/queries"
	"mall/stores/internal/domain"
)

type Application struct {
	application.App
	logger *slog.Logger
}

var _ application.App = (*Application)(nil)

func LogApplicationAccess(application application.App, logger *slog.Logger) Application {
	return Application{
		App:    application,
		logger: logger,
	}
}

func (a Application) CreateStore(ctx context.Context, cmd commands.CreateStoreRequest) (err error) {
	a.logger.Info("--> Stores.CreateStore")
	defer func() {
		if err != nil {
			a.logger.Error(err.Error())
		}
		a.logger.Info("<-- Stores.CreateStore")
	}()

	return a.App.CreateStore(ctx, cmd)
}

func (a Application) EnableParticipation(ctx context.Context, cmd commands.EnableParticipationRequest) (err error) {
	a.logger.Info("--> Stores.EnableParticipation")
	defer func() {
		if err != nil {
			a.logger.Error(err.Error())
		}
		a.logger.Info("<-- Stores.EnableParticipation")
	}()

	return a.App.EnableParticipation(ctx, cmd)
}

func (a Application) DisableParticipation(ctx context.Context, cmd commands.DisableParticipationRequest) (err error) {
	a.logger.Info("--> Stores.DisableParticipation")
	defer func() {
		if err != nil {
			a.logger.Error(err.Error())
		}
		a.logger.Info("<-- Stores.DisableParticipation")
	}()

	return a.App.DisableParticipation(ctx, cmd)
}

func (a Application) AddProduct(ctx context.Context, cmd commands.AddProductRequest) (err error) {
	a.logger.Info("--> Stores.AddProduct")
	defer func() {
		if err != nil {
			a.logger.Error(err.Error())
		}
		a.logger.Info("<-- Stores.AddProduct")
	}()

	return a.App.AddProduct(ctx, cmd)
}

func (a Application) RemoveProduct(ctx context.Context, cmd commands.RemoveProductRequest) (err error) {
	a.logger.Info("--> Stores.RemoveProduct")
	defer func() {
		if err != nil {
			a.logger.Error(err.Error())
		}
		a.logger.Info("<-- Stores.RemoveProduct")
	}()

	return a.App.RemoveProduct(ctx, cmd)
}

func (a Application) GetStore(ctx context.Context, query queries.GetStoreRequest) (store *domain.MallStore, err error) {
	a.logger.Info("--> Stores.GetStore")
	defer func() {
		if err != nil {
			a.logger.Error(err.Error())
		}
		a.logger.Info("<-- Stores.GetStore")
	}()

	return a.App.GetStore(ctx, query)
}

func (a Application) GetStores(ctx context.Context, query queries.GetStoresRequest) (stores []*domain.MallStore, err error) {
	a.logger.Info("--> Stores.GetStores")
	defer func() {
		if err != nil {
			a.logger.Error(err.Error())
		}
		a.logger.Info("<-- Stores.GetStores")
	}()

	return a.App.GetStores(ctx, query)
}

func (a Application) GetParticipatingStores(ctx context.Context, query queries.GetParticipatingStoreRequest) (stores []*domain.MallStore, err error) {
	a.logger.Info("--> Stores.GetParticipatingStores")
	defer func() {
		if err != nil {
			a.logger.Error(err.Error())
		}
		a.logger.Info("<-- Stores.GetParticipatingStores")
	}()

	return a.App.GetParticipatingStores(ctx, query)
}

func (a Application) GetCatalog(ctx context.Context, query queries.GetCatalogRequest) (products []*domain.CatalogProduct, err error) {
	a.logger.Info("--> Stores.GetCatalog")
	defer func() {
		if err != nil {
			a.logger.Error(err.Error())
		}
		a.logger.Info("<-- Stores.GetCatalog")
	}()

	return a.App.GetCatalog(ctx, query)
}

func (a Application) GetProduct(ctx context.Context, query queries.GetProductRequest) (product *domain.CatalogProduct, err error) {
	a.logger.Info("--> Stores.GetProduct")
	defer func() {
		if err != nil {
			a.logger.Error(err.Error())
		}
		a.logger.Info("<-- Stores.GetProduct")
	}()

	return a.App.GetProduct(ctx, query)
}
