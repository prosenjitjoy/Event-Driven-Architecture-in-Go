package domain

import (
	"errors"
)

var (
	ErrProductNameIsBlank     = errors.New("product name cannot be blank")
	ErrProductPriceIsNegative = errors.New("product price cannot be negative")
)

type Product struct {
	ID          string
	StoreID     string
	Name        string
	Description string
	SKU         string
	Price       float64
}

func CreateProduct(id, storeID, name, description, sku string, price float64) (*Product, error) {
	if name == "" {
		return nil, ErrProductNameIsBlank
	}

	if price < 0 {
		return nil, ErrProductPriceIsNegative
	}

	product := &Product{
		ID:          id,
		StoreID:     storeID,
		Name:        name,
		Description: description,
		SKU:         sku,
		Price:       price,
	}

	return product, nil
}
