package commands

import (
	"context"
	"mall/internal/ddd"
	"mall/stores/internal/domain"
)

type RebrandProductRequest struct {
	ID          string
	Name        string
	Description string
}

type RebrandProductHandler struct {
	products  domain.ProductRepository
	publisher ddd.EventPublisher[ddd.Event]
}

func NewRebrandProductHandler(products domain.ProductRepository, publisher ddd.EventPublisher[ddd.Event]) RebrandProductHandler {
	return RebrandProductHandler{
		products:  products,
		publisher: publisher,
	}
}

func (h RebrandProductHandler) RebrandProduct(ctx context.Context, cmd RebrandProductRequest) error {
	product, err := h.products.Load(ctx, cmd.ID)
	if err != nil {
		return err
	}

	event, err := product.Rebrand(cmd.Name, cmd.Description)
	if err != nil {
		return err
	}

	if err := h.products.Save(ctx, product); err != nil {
		return err
	}

	if err := h.publisher.Publish(ctx, event); err != nil {
		return err
	}

	return nil
}
