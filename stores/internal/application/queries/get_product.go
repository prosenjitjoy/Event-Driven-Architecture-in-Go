package queries

import (
	"context"
	"mall/stores/internal/domain"
)

type GetProductRequest struct {
	ID string
}

type GetProductHandler struct {
	products domain.ProductRepository
}

func NewGetProductHandler(products domain.ProductRepository) GetProductHandler {
	return GetProductHandler{
		products: products,
	}
}

func (h GetProductHandler) GetProduct(ctx context.Context, query GetProductRequest) (*domain.Product, error) {
	return h.products.Find(ctx, query.ID)
}
