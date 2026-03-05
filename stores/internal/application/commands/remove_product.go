package commands

import (
	"context"
	"mall/internal/ddd"
	"mall/stores/internal/domain"
)

type RemoveProductRequest struct {
	ID string
}

type RemoveProductHandler struct {
	products  domain.ProductRepository
	publisher ddd.EventPublisher[ddd.Event]
}

func NewRemoveProductHandler(products domain.ProductRepository, publisher ddd.EventPublisher[ddd.Event]) RemoveProductHandler {
	return RemoveProductHandler{
		products:  products,
		publisher: publisher,
	}
}

func (h RemoveProductHandler) RemoveProduct(ctx context.Context, cmd RemoveProductRequest) error {
	product, err := h.products.Load(ctx, cmd.ID)
	if err != nil {
		return err
	}

	event, err := product.Remove()
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
