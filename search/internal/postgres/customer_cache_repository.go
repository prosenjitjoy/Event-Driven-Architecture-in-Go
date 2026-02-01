package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"mall/internal/postgres"
	"mall/search/internal/domain"

	"github.com/lib/pq"
)

type CustomerCacheRepository struct {
	tableName string
	db        postgres.DBTX
	fallback  domain.CustomerRepository
}

var _ domain.CustomerCacheRepository = (*CustomerCacheRepository)(nil)

func NewCustomerCacheRepository(tableName string, db postgres.DBTX, fallback domain.CustomerRepository) CustomerCacheRepository {
	return CustomerCacheRepository{
		tableName: tableName,
		db:        db,
		fallback:  fallback,
	}
}

func (r CustomerCacheRepository) Add(ctx context.Context, customerID, name string) error {
	const query = "INSERT INTO %s (id, name) VALUES($1, $2)"

	_, err := r.db.ExecContext(ctx, r.table(query), customerID, name)
	if err != nil {
		if pgErr, ok := err.(*pq.Error); ok {
			// unique_violation error
			if pgErr.Code == pq.ErrorCode("23505") {
				return nil
			}
		}
		return err
	}

	return nil
}

func (r CustomerCacheRepository) Find(ctx context.Context, customerID string) (*domain.Customer, error) {
	const query = "SELECT name FROM %s WHERE id = $1 LIMIT 1"

	customer := &domain.Customer{
		ID: customerID,
	}

	err := r.db.QueryRowContext(ctx, r.table(query), customerID).Scan(&customer.Name)
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("scanning customer: %w", err)
		}

		customer, err := r.fallback.Find(ctx, customerID)
		if err != nil {
			return nil, fmt.Errorf("customer fallback failed: %w", err)
		}

		// attempt to add it to the cache
		return customer, r.Add(ctx, customer.ID, customer.Name)
	}

	return customer, nil
}

func (r CustomerCacheRepository) table(query string) string {
	return fmt.Sprintf(query, r.tableName)
}
