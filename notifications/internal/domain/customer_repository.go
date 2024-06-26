package domain

import (
	"context"
)

type CustomerRepository interface {
	Find(ctx context.Context, customerID string) (*Customer, error)
}
