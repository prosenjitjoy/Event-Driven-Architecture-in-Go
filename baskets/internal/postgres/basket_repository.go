package postgres

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"mall/baskets/internal/domain"
	"mall/internal/ddd"
)

type BasketRepository struct {
	tableName string
	db        *sql.DB
}

var _ domain.BasketRepository = (*BasketRepository)(nil)

func NewBasketRepository(tableName string, db *sql.DB) BasketRepository {
	return BasketRepository{
		tableName: tableName,
		db:        db,
	}
}

func (r BasketRepository) Find(ctx context.Context, basketID string) (*domain.Basket, error) {
	const query = "SELECT customer_id, payment_id, items, status FROM %s WHERE id = $1 LIMIT 1"

	basket := &domain.Basket{
		AggregateBase: ddd.AggregateBase{ID: basketID},
	}
	var items []byte
	var status string

	err := r.db.QueryRowContext(ctx, r.table(query), basketID).Scan(&basket.CustomerID, &basket.PaymentID, &items, &status)
	if err != nil {
		return nil, fmt.Errorf("INTERNAL_SERVER_ERROR: %w", err)
	}

	basket.Status = domain.ToBasketStatus(status)

	if err = json.Unmarshal(items, &basket.Items); err != nil {
		return nil, fmt.Errorf("INTERNAL_SERVER_ERROR: %w", err)
	}

	return basket, nil
}

func (r BasketRepository) Save(ctx context.Context, basket *domain.Basket) error {
	const query = "INSERT INTO %s (id, customer_id, payment_id, items, status) VALUES ($1, $2, $3, $4, $5)"

	items, err := json.Marshal(basket.Items)
	if err != nil {
		return fmt.Errorf("INTERNAL_SERVER_ERROR: %w", err)
	}

	_, err = r.db.ExecContext(ctx, r.table(query), basket.ID, basket.CustomerID, basket.PaymentID, items, basket.Status.String())
	if err != nil {
		return fmt.Errorf("INTERNAL_SERVER_ERROR: %w", err)
	}

	return nil
}

func (r BasketRepository) Update(ctx context.Context, basket *domain.Basket) error {
	const query = "UPDATE %s SET customer_id = $2, payment_id = $3, items = $4, status = $5 WHERE id = $1"

	items, err := json.Marshal(basket.Items)
	if err != nil {
		return fmt.Errorf("INTERNAL_SERVER_ERROR: %w", err)
	}

	_, err = r.db.ExecContext(ctx, r.table(query), basket.ID, basket.CustomerID, basket.PaymentID, items, basket.Status.String())
	if err != nil {
		return fmt.Errorf("INTERNAL_SERVER_ERROR: %w", err)
	}

	return nil
}

func (r BasketRepository) DeleteBasket(ctx context.Context, basketID string) error {
	const query = "DELETE FROM %s WHERE id = $1 LIMIT 1"
	_, err := r.db.ExecContext(ctx, r.table(query), basketID)
	if err != nil {
		return fmt.Errorf("INTERNAL_SERVER_ERROR: %w", err)
	}

	return nil
}

func (r BasketRepository) table(query string) string {
	return fmt.Sprintf(query, r.tableName)
}
