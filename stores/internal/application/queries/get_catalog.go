package queries

import (
	"context"
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
	return h.catalog.GetCatalog(ctx, query.StoreID)
}
