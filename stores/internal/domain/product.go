package domain

import (
	"errors"
	"fmt"
	"mall/internal/ddd"
	"mall/internal/es"
)

const ProductAggregate = "stores.Product"

var (
	ErrProductNameIsBlank     = errors.New("the product name cannot be blank")
	ErrProductPriceIsNegative = errors.New("the product price cannot be negative")
	ErrNotAPriceIncrease      = errors.New("the price change would be a decrease")
	ErrNotAPriceDecrease      = errors.New("the price change would be an increase")
)

type Product struct {
	es.Aggregate
	StoreID     string
	Name        string
	Description string
	SKU         string
	Price       float64
}

var _ interface {
	es.EventApplier
	es.Snapshotter
} = (*Product)(nil)

func NewProduct(id string) *Product {
	return &Product{
		Aggregate: es.NewAggregate(id, ProductAggregate),
	}
}

func CreateProduct(id, storeID, name, description, sku string, price float64) (*Product, error) {
	if name == "" {
		return nil, ErrProductNameIsBlank
	}

	if price < 0 {
		return nil, ErrProductPriceIsNegative
	}

	product := NewProduct(id)

	product.AddEvent(ProductAddedEvent, &ProductAdded{
		StoreID:     storeID,
		Name:        name,
		Description: description,
		SKU:         sku,
		Price:       price,
	})

	return product, nil
}

func (p *Product) Rebrand(name, description string) error {
	p.AddEvent(ProductRebrandedEvent, &ProductRebranded{
		Name:        name,
		Description: description,
	})

	return nil
}

func (p *Product) IncreasePrice(price float64) error {
	if price < p.Price {
		return ErrNotAPriceIncrease
	}

	p.AddEvent(ProductPriceIncreasedEvent, &ProductPriceChanged{
		Delta: price - p.Price,
	})

	return nil
}

func (p *Product) DecreasePrice(price float64) error {
	if price > p.Price {
		return ErrNotAPriceDecrease
	}

	p.AddEvent(ProductPriceDecreasedEvent, &ProductPriceChanged{
		Delta: price - p.Price,
	})

	return nil
}

func (p *Product) Remove() error {
	p.AddEvent(ProductRemovedEvent, &ProductRemoved{})

	return nil
}

func (p *Product) ApplyEvent(event ddd.Event) error {
	switch payload := event.Payload().(type) {
	case *ProductAdded:
		p.StoreID = payload.StoreID
		p.Name = payload.Name
		p.Description = payload.Description
		p.SKU = payload.SKU
		p.Price = payload.Price
	case *ProductRebranded:
		p.Name = payload.Name
		p.Description = payload.Description
	case *ProductPriceChanged:
		p.Price = p.Price + payload.Delta
	case *ProductRemoved:
		// no operation
	default:
		return fmt.Errorf("%T received the event %s with unexpected payload %T", p, event.EventName(), payload)
	}

	return nil
}

func (p *Product) ApplySnapshot(snapshot es.Snapshot) error {
	switch ss := snapshot.(type) {
	case *ProductV1:
		p.StoreID = ss.StoreID
		p.Name = ss.Name
		p.Description = ss.Description
		p.SKU = ss.SKU
		p.Price = ss.Price
	default:
		return fmt.Errorf("%T received the unexpected snapshot %T", p, snapshot)
	}

	return nil
}

func (p Product) ToSnapshot() es.Snapshot {
	return ProductV1{
		StoreID:     p.StoreID,
		Name:        p.Name,
		Description: p.Description,
		SKU:         p.SKU,
		Price:       p.Price,
	}
}
