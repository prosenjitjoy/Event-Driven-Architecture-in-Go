package application

import (
	"context"
	"fmt"
	"mall/payments/internal/domain"
)

type AuthorizePayment struct {
	ID         string
	CustomerID string
	Amount     float64
}

type ConfirmPayment struct {
	ID string
}

type CreateInvoice struct {
	ID        string
	OrderID   string
	PaymentID string
	Amount    float64
}

type AdjustInvoice struct {
	ID     string
	Amount float64
}

type PayInvoice struct {
	ID string
}

type CancelInvoice struct {
	ID string
}

type App interface {
	AuthorizePayment(ctx context.Context, authorize AuthorizePayment) error
	ConfirmPayment(ctx context.Context, confirm ConfirmPayment) error
	CreateInvoice(ctx context.Context, create CreateInvoice) error
	AdjustInvoice(ctx context.Context, adjust AdjustInvoice) error
	PayInvoice(ctx context.Context, pay PayInvoice) error
	CancelInvoice(ctx context.Context, cancel CancelInvoice) error
}

type Application struct {
	invoices domain.InvoiceRepository
	payments domain.PaymentRepository
	orders   domain.OrderRepsitory
}

var _ App = (*Application)(nil)

func New(invoices domain.InvoiceRepository, payments domain.PaymentRepository, orders domain.OrderRepsitory) *Application {
	return &Application{
		invoices: invoices,
		payments: payments,
		orders:   orders,
	}
}

func (a Application) AuthorizePayment(ctx context.Context, authorize AuthorizePayment) error {
	return a.payments.Save(ctx, &domain.Payment{
		ID:         authorize.ID,
		CustomerID: authorize.CustomerID,
		Amount:     authorize.Amount,
	})
}

func (a Application) ConfirmPayment(ctx context.Context, confirm ConfirmPayment) error {
	_, err := a.payments.Find(ctx, confirm.ID)
	if err != nil {
		return fmt.Errorf("NOT_FOUNT: %w", err)
	}

	return nil
}

func (a Application) CreateInvoice(ctx context.Context, create CreateInvoice) error {
	return a.invoices.Save(ctx, &domain.Invoice{
		ID:      create.ID,
		OrderID: create.OrderID,
		Amount:  create.Amount,
		Status:  domain.InvoicePending,
	})
}

func (a Application) AdjustInvoice(ctx context.Context, adjust AdjustInvoice) error {
	invoice, err := a.invoices.Find(ctx, adjust.ID)
	if err != nil {
		return err
	}

	invoice.Amount = adjust.Amount

	return a.invoices.Update(ctx, invoice)
}

func (a Application) PayInvoice(ctx context.Context, pay PayInvoice) error {
	invoice, err := a.invoices.Find(ctx, pay.ID)
	if err != nil {
		return err
	}

	if invoice.Status != domain.InvoicePending {
		return fmt.Errorf("BAD_REQUEST: %s", "invoice cannot be paid for")
	}

	invoice.Status = domain.InvoicePaid

	err = a.orders.Complete(ctx, invoice.ID, invoice.OrderID)
	if err != nil {
		return err
	}

	return a.invoices.Update(ctx, invoice)
}

func (a Application) CancelInvoice(ctx context.Context, cancel CancelInvoice) error {
	invoice, err := a.invoices.Find(ctx, cancel.ID)
	if err != nil {
		return err
	}

	if invoice.Status != domain.InvoicePending {
		return fmt.Errorf("BAD_REQUEST: %s", "invoice cannot be paid for")
	}

	invoice.Status = domain.InvoiceCanceled

	return a.invoices.Update(ctx, invoice)
}
