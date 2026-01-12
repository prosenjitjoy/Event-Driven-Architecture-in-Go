package grpc

import (
	"context"
	"mall/ordering/internal/domain"
	"mall/payments/paymentspb"

	"google.golang.org/grpc"
)

type InvoiceRepository struct {
	client paymentspb.PaymentsServiceClient
}

var _ domain.InvoiceRepository = (*InvoiceRepository)(nil)

func NewInvoiceRepository(conn *grpc.ClientConn) InvoiceRepository {
	return InvoiceRepository{
		client: paymentspb.NewPaymentsServiceClient(conn),
	}
}

func (r InvoiceRepository) Save(ctx context.Context, orderID, paymentID string, amount float64) error {
	_, err := r.client.CreateInvoice(ctx, &paymentspb.CreateInvoiceRequest{
		OrderId:   orderID,
		PaymentId: paymentID,
		Amount:    amount,
	})
	if err != nil {
		return err
	}

	return nil
}

func (r InvoiceRepository) Delete(ctx context.Context, invoiceID string) error {
	_, err := r.client.CancelInvoice(ctx, &paymentspb.CancelInvoiceRequest{
		Id: invoiceID,
	})
	if err != nil {
		return err
	}

	return nil
}
