package queries

import (
	"context"
	"mall/stores/internal/domain"
)

type GetProductRequest struct {
	ID string
}

type GetProductHandler struct {
	catalog domain.CatalogRepository
}

func NewGetProductHandler(catalog domain.CatalogRepository) GetProductHandler {
	return GetProductHandler{
		catalog: catalog,
	}
}

func (h GetProductHandler) GetProduct(ctx context.Context, query GetProductRequest) (*domain.CatalogProduct, error) {
	return h.catalog.Find(ctx, query.ID)
}
