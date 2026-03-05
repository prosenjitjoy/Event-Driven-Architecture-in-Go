package commands

import (
	"context"
	"mall/internal/ddd"
	"mall/stores/internal/domain"
)

type DecreaseProductPriceRequest struct {
	ID    string
	Price float64
}

type DecreaseProductPriceHandler struct {
	products  domain.ProductRepository
	publisher ddd.EventPublisher[ddd.Event]
}

func NewDecreaseProductPriceHandler(products domain.ProductRepository, publisher ddd.EventPublisher[ddd.Event]) DecreaseProductPriceHandler {
	return DecreaseProductPriceHandler{
		products:  products,
		publisher: publisher,
	}
}

func (h DecreaseProductPriceHandler) DecreaseProductPrice(ctx context.Context, cmd DecreaseProductPriceRequest) error {
	product, err := h.products.Load(ctx, cmd.ID)
	if err != nil {
		return err
	}

	event, err := product.DecreasePrice(cmd.Price)
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
