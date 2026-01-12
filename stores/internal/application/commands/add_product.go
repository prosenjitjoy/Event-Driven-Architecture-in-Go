package commands

import (
	"context"
	"fmt"
	"mall/internal/ddd"
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
	stores          domain.StoreRepository
	products        domain.ProductRepository
	domainPublisher ddd.EventPublisher
}

func NewAddProductHandler(stores domain.StoreRepository, products domain.ProductRepository, domainPublisher ddd.EventPublisher) AddProductHandler {
	return AddProductHandler{
		stores:          stores,
		products:        products,
		domainPublisher: domainPublisher,
	}
}

func (h AddProductHandler) AddProduct(ctx context.Context, cmd AddProductRequest) error {
	_, err := h.stores.Find(ctx, cmd.StoreID)
	if err != nil {
		return fmt.Errorf("error adding product: %w", err)
	}

	product, err := domain.CreateProduct(cmd.ID, cmd.StoreID, cmd.Name, cmd.Description, cmd.SKU, cmd.Price)
	if err != nil {
		return fmt.Errorf("error adding product: %w", err)
	}

	if err = h.products.Save(ctx, product); err != nil {
		return fmt.Errorf("error adding product: %w", err)
	}

	if err = h.domainPublisher.Publish(ctx, product.GetEvents()...); err != nil {
		return err
	}

	return nil
}
