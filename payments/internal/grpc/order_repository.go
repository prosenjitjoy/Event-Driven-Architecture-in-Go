package grpc

import (
	"context"
	"mall/ordering/orderingpb"
	"mall/payments/internal/domain"

	"google.golang.org/grpc"
)

type OrderRepsitory struct {
	client orderingpb.OrderingServiceClient
}

var _ domain.OrderRepsitory = (*OrderRepsitory)(nil)

func NewOrderRepository(conn *grpc.ClientConn) OrderRepsitory {
	return OrderRepsitory{
		client: orderingpb.NewOrderingServiceClient(conn),
	}
}

func (r OrderRepsitory) Complete(ctx context.Context, invoiceID, orderID string) error {
	_, err := r.client.CompleteOrder(ctx, &orderingpb.CompleteOrderRequest{
		Id:        orderID,
		InvoiceId: invoiceID,
	})
	if err != nil {
		return err
	}

	return nil
}
