package commands

import (
	"context"
	"fmt"
	"mall/stores/internal/domain"
)

type AddProductRequest struct {
	ID          string
	StoreID     string
	Name        string
	Description string
	SKU         string
	Price       float64
}

type AddProductHandler struct {
	products domain.ProductRepository
}

func NewAddProductHandler(products domain.ProductRepository) AddProductHandler {
	return AddProductHandler{
		products: products,
	}
}

func (h AddProductHandler) AddProduct(ctx context.Context, cmd AddProductRequest) error {
	product, err := domain.CreateProduct(cmd.ID, cmd.StoreID, cmd.Name, cmd.Description, cmd.SKU, cmd.Price)
	if err != nil {
		return fmt.Errorf("error adding product: %w", err)
	}

	if err = h.products.Save(ctx, product); err != nil {
		return fmt.Errorf("error adding product: %w", err)
	}

	return nil
}
