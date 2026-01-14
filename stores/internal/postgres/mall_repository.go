package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"mall/stores/internal/domain"
)

type MallRepository struct {
	tableName string
	db        *sql.DB
}

var _ domain.MallRepository = (*MallRepository)(nil)

func NewMallRepository(tableName string, db *sql.DB) MallRepository {
	return MallRepository{
		tableName: tableName,
		db:        db,
	}
}

func (r MallRepository) AddStore(ctx context.Context, storeID, name, location string) error {
	const query = "INSERT INTO %s (id, name, location, participating) VALUES ($1, $2, $3, $4)"

	_, err := r.db.ExecContext(ctx, r.table(query), storeID, name, location, false)
	if err != nil {
		return err
	}

	return nil
}

func (r MallRepository) SetStoreParticipation(ctx context.Context, storeID string, participating bool) error {
	const query = "UPDATE %s SET participating = $2 WHERE id = $1"

	_, err := r.db.ExecContext(ctx, r.table(query), storeID, participating)
	if err != nil {
		return err
	}

	return nil
}

func (r MallRepository) RenameStore(ctx context.Context, storeID, name string) error {
	const query = "UPDATE %s SET name = $2 WHERE id = $1"

	_, err := r.db.ExecContext(ctx, r.table(query), storeID, name)
	if err != nil {
		return err
	}

	return nil
}

func (r MallRepository) Find(ctx context.Context, storeID string) (*domain.MallStore, error) {
	const query = "SELECT name, location, participating FROM %s WHERE id = $1 LIMIT 1"

	store := &domain.MallStore{
		ID: storeID,
	}

	err := r.db.QueryRowContext(ctx, r.table(query), storeID).Scan(&store.Name, &store.Location, &store.Participating)
	if err != nil {
		return nil, fmt.Errorf("scanning store: %w", err)
	}

	return store, nil
}

func (r MallRepository) All(ctx context.Context) ([]*domain.MallStore, error) {
	const query = "SELECT id, name, location, participating FROM %s"

	rows, err := r.db.QueryContext(ctx, r.table(query))
	if err != nil {
		return nil, fmt.Errorf("querying stores: %w", err)
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			err = fmt.Errorf("closing store rows: %w", err)
		}
	}(rows)

	stores := make([]*domain.MallStore, 0)
	for rows.Next() {
		store := &domain.MallStore{}
		err := rows.Scan(&store.ID, &store.Name, &store.Location, &store.Participating)
		if err != nil {
			return nil, fmt.Errorf("scanning store: %w", err)
		}

		stores = append(stores, store)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("finishing store rows: %w", err)
	}

	return stores, nil
}

func (r MallRepository) AllParticipating(ctx context.Context) ([]*domain.MallStore, error) {
	const query = "SELECT id, name, location, participating FROM %s WHERE participating is true"

	rows, err := r.db.QueryContext(ctx, r.table(query))
	if err != nil {
		return nil, fmt.Errorf("querying participating stores: %w", err)
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			err = fmt.Errorf("closing participating store rows: %w", err)
		}
	}(rows)

	stores := make([]*domain.MallStore, 0)
	for rows.Next() {
		store := &domain.MallStore{}
		err := rows.Scan(&store.ID, &store.Name, &store.Location, &store.Participating)
		if err != nil {
			return nil, fmt.Errorf("scanning store: %w", err)
		}

		stores = append(stores, store)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("finishing participating store rows: %w", err)
	}

	return stores, nil
}

func (r MallRepository) table(query string) string {
	return fmt.Sprintf(query, r.tableName)
}
