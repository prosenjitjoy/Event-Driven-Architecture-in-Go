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
	products  domain.ProductRepository
	publisher ddd.EventPublisher[ddd.Event]
}

func NewAddProductHandler(products domain.ProductRepository, publisher ddd.EventPublisher[ddd.Event]) AddProductHandler {
	return AddProductHandler{
		products:  products,
		publisher: publisher,
	}
}

func (h AddProductHandler) AddProduct(ctx context.Context, cmd AddProductRequest) error {
	product, err := h.products.Load(ctx, cmd.ID)
	if err != nil {
		return fmt.Errorf("error adding product: %w", err)
	}

	event, err := product.InitProduct(cmd.StoreID, cmd.Name, cmd.Description, cmd.SKU, cmd.Price)
	if err != nil {
		return fmt.Errorf("initializing product: %w", err)
	}

	if err := h.products.Save(ctx, product); err != nil {
		return fmt.Errorf("error adding product: %w", err)
	}

	if err := h.publisher.Publish(ctx, event); err != nil {
		return fmt.Errorf("publishing domain event: %w", err)
	}

	return nil
}
