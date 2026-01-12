package queries

import (
	"context"
	"mall/stores/internal/domain"
)

type GetCatalogRequest struct {
	StoreID string
}

type GetCatalogHandler struct {
	products domain.ProductRepository
}

func NewGetCatalogHandler(product domain.ProductRepository) GetCatalogHandler {
	return GetCatalogHandler{products: product}
}

func (h GetCatalogHandler) GetCatalog(ctx context.Context, query GetCatalogRequest) ([]*domain.Product, error) {
	return h.products.GetCatalog(ctx, query.StoreID)
}
