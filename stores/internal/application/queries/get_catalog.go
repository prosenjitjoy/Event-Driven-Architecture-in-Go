package queries

import (
	"context"
	"log/slog"
	"mall/stores/internal/domain"
)

type GetCatalogRequest struct {
	StoreID string
}

type GetCatalogHandler struct {
	catalog domain.CatalogRepository
}

func NewGetCatalogHandler(catalog domain.CatalogRepository) GetCatalogHandler {
	return GetCatalogHandler{
		catalog: catalog,
	}
}

func (h GetCatalogHandler) GetCatalog(ctx context.Context, query GetCatalogRequest) ([]*domain.CatalogProduct, error) {
	slog.Info("--> here")
	defer slog.Info("<-- here")
	return h.catalog.GetCatalog(ctx, query.StoreID)
}
