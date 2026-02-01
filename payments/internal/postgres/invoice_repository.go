package postgres

import (
	"context"
	"fmt"
	"mall/internal/postgres"
	"mall/payments/internal/domain"
)

type InvoiceRepository struct {
	tableName string
	db        postgres.DBTX
}

var _ domain.InvoiceRepository = (*InvoiceRepository)(nil)

func NewInvoiceRepository(tableName string, db postgres.DBTX) InvoiceRepository {
	return InvoiceRepository{
		tableName: tableName,
		db:        db,
	}
}

func (r InvoiceRepository) Find(ctx context.Context, invoiceID string) (*domain.Invoice, error) {
	const query = "SELECT order_id, amount, status FROM %s WHERE id = $1 LIMIT 1"

	invoice := &domain.Invoice{
		ID: invoiceID,
	}
	var status string

	err := r.db.QueryRowContext(ctx, r.table(query), invoiceID).Scan(&invoice.OrderID, &invoice.Amount, &status)
	if err != nil {
		return nil, fmt.Errorf("scanning invoice: %w", err)
	}

	invoice.Status, err = r.statusToDomain(status)
	if err != nil {
		return nil, err
	}

	return invoice, nil
}

func (r InvoiceRepository) Save(ctx context.Context, invoice *domain.Invoice) error {
	const query = "INSERT INTO %s (id, order_id, amount, status) VALUES ($1, $2, $3, $4)"

	_, err := r.db.ExecContext(ctx, r.table(query), invoice.ID, invoice.OrderID, invoice.Amount, invoice.Status.String())
	if err != nil {
		return err
	}

	return nil
}

func (r InvoiceRepository) Update(ctx context.Context, invoice *domain.Invoice) error {
	const query = "UPDATE %s SET amount = $2, status = $3 WHERE id = $1"

	_, err := r.db.ExecContext(ctx, r.table(query), invoice.ID, invoice.Amount, invoice.Status.String())
	if err != nil {
		return err
	}

	return nil
}

func (r InvoiceRepository) table(query string) string {
	return fmt.Sprintf(query, r.tableName)
}

func (r InvoiceRepository) statusToDomain(status string) (domain.InvoiceStatus, error) {
	switch status {
	case domain.InvoiceIsPending.String():
		return domain.InvoiceIsPending, nil
	case domain.InvoiceIsPaid.String():
		return domain.InvoiceIsPaid, nil
	default:
		return domain.InvoiceIsUnknown, fmt.Errorf("unknown invoice status: %s", status)
	}
}
