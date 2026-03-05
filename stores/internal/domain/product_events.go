package domain

type ProductAdded struct {
	StoreID     string
	Name        string
	Description string
	SKU         string
	Price       float64
}

type ProductRebranded struct {
	Name        string
	Description string
}

type ProductPriceChanged struct {
	Delta float64
}

type ProductRemoved struct{}

// domain events
type ProductPriceDelta struct {
	Product *Product
	Delta   float64
}
