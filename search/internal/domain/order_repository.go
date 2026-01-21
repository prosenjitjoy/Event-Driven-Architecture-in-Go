package domain

import (
	"context"
	"time"
)

type Filters struct {
	CustomerID string
	After      time.Time
	Before     time.Time
	StoreIDs   []string
	ProductIDs []string
	MinTotal   float64
	MaxTotal   float64
	Status     string
}

type SearchOrders struct {
	Filters Filters
	Next    string
	Limit   int
}

type GetOrder struct {
	OrderID string
}

type OrderRepository interface {
	Add(ctx context.Context, order *Order) error
	UpdateStatus(ctx context.Context, orderID, status string) error
	Search(ctx context.Context, search SearchOrders) ([]Order, error)
	Get(ctx context.Context, orderID string) (*Order, error)
}
