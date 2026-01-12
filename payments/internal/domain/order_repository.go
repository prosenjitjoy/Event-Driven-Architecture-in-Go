package domain

import "context"

type OrderRepsitory interface {
	Complete(ctx context.Context, invoiceID, orderID string) error
}
