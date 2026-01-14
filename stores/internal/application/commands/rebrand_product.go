package commands

import (
	"context"
	"mall/stores/internal/domain"
)

type RebrandProductRequest struct {
	ID          string
	Name        string
	Description string
}

type RebrandProductHandler struct {
	products domain.ProductRepository
}

func NewRebrandProductHandler(products domain.ProductRepository) RebrandProductHandler {
	return RebrandProductHandler{
		products: products,
	}
}

func (h RebrandProductHandler) RebrandProduct(ctx context.Context, cmd RebrandProductRequest) error {
	product, err := h.products.Load(ctx, cmd.ID)
	if err != nil {
		return err
	}

	if err = product.Rebrand(cmd.Name, cmd.Description); err != nil {
		return err
	}

	return h.products.Save(ctx, product)
}
