package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"mall/search/internal/domain"

	"github.com/lib/pq"
)

type StoreCacheRepository struct {
	tableName string
	db        *sql.DB
	fallback  domain.StoreRepository
}

var _ domain.StoreCacheRepository = (*StoreCacheRepository)(nil)

func NewStoreCacheRepository(tableName string, db *sql.DB, fallback domain.StoreRepository) StoreCacheRepository {
	return StoreCacheRepository{
		tableName: tableName,
		db:        db,
		fallback:  fallback,
	}
}

func (r StoreCacheRepository) Add(ctx context.Context, storeID, name string) error {
	const query = "INSERT INTO %s (id, name) VALUES ($1, $2)"

	_, err := r.db.ExecContext(ctx, r.table(query), storeID, name)
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

func (r StoreCacheRepository) Rename(ctx context.Context, storeID, name string) error {
	const query = "UPDATE %s SET name = $2 WHERE id = $1"

	_, err := r.db.ExecContext(ctx, r.table(query), storeID, name)
	if err != nil {
		return err
	}

	return nil
}

func (r StoreCacheRepository) Find(ctx context.Context, storeID string) (*domain.Store, error) {
	const query = "SELECT name FROM %s WHERE id = $1 LIMIT 1"

	store := &domain.Store{
		ID: storeID,
	}

	err := r.db.QueryRowContext(ctx, r.table(query), storeID).Scan(&store.Name)
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("scanning store: %w", err)
		}

		store, err := r.fallback.Find(ctx, storeID)
		if err != nil {
			return nil, fmt.Errorf("store fallback failed: %w", err)
		}

		// attempt to add it to the cache
		return store, r.Add(ctx, store.ID, store.Name)
	}

	return store, nil
}

func (r StoreCacheRepository) table(query string) string {
	return fmt.Sprintf(query, r.tableName)
}
