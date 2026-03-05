package commands

import (
	"context"
	"mall/internal/ddd"
	"mall/stores/internal/domain"
)

type IncreaseProductPriceRequest struct {
	ID    string
	Price float64
}

type IncreaseProductPriceHandler struct {
	products  domain.ProductRepository
	publisher ddd.EventPublisher[ddd.Event]
}

func NewIncreaseProductPriceHandler(products domain.ProductRepository, publisher ddd.EventPublisher[ddd.Event]) IncreaseProductPriceHandler {
	return IncreaseProductPriceHandler{
		products:  products,
		publisher: publisher,
	}
}

func (h IncreaseProductPriceHandler) IncreaseProductPrice(ctx context.Context, cmd IncreaseProductPriceRequest) error {
	product, err := h.products.Load(ctx, cmd.ID)
	if err != nil {
		return err
	}

	event, err := product.IncreasePrice(cmd.Price)
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
