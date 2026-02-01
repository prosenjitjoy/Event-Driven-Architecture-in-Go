package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"mall/internal/postgres"
	"mall/stores/internal/domain"
)

type CatalogRepository struct {
	tableName string
	db        postgres.DBTX
}

var _ domain.CatalogRepository = (*CatalogRepository)(nil)

func NewCatalogRepository(tableName string, db postgres.DBTX) CatalogRepository {
	return CatalogRepository{
		tableName: tableName,
		db:        db,
	}
}

func (r CatalogRepository) AddProduct(ctx context.Context, productID, storeID, name, description, sku string, price float64) error {
	const query = "INSERT INTO %s (id, store_id, name, description, sku, price) VALUES ($1, $2, $3, $4, $5, $6)"

	_, err := r.db.ExecContext(ctx, r.table(query), productID, storeID, name, description, sku, price)
	if err != nil {
		return err
	}

	return nil
}

func (r CatalogRepository) Rebrand(ctx context.Context, productID, name, description string) error {
	const query = "UPDATE %s SET name = $2, description = $3 WHERE id = $1"

	_, err := r.db.ExecContext(ctx, r.table(query), productID, name, description)
	if err != nil {
		return err
	}

	return nil
}

func (r CatalogRepository) UpdatePrice(ctx context.Context, productID string, delta float64) error {
	const query = "UPDATE %s SET price = price + $2 WHERE id = $1"

	_, err := r.db.ExecContext(ctx, r.table(query), productID, delta)
	if err != nil {
		return err
	}

	return nil
}

func (r CatalogRepository) RemoveProduct(ctx context.Context, productID string) error {
	const query = "DELETE FROM %s WHERE id = $1"

	_, err := r.db.ExecContext(ctx, r.table(query), productID)
	if err != nil {
		return err
	}

	return nil
}

func (r CatalogRepository) Find(ctx context.Context, productID string) (*domain.CatalogProduct, error) {
	const query = "SELECT store_id, name, description, sku, price FROM %s WHERE id = $1 LIMIT 1"

	product := &domain.CatalogProduct{
		ID: productID,
	}

	err := r.db.QueryRowContext(ctx, r.table(query), productID).Scan(&product.StoreID, &product.Name, &product.Description, &product.SKU, &product.Price)
	if err != nil {
		return nil, fmt.Errorf("scanning product: %w", err)
	}

	return product, nil
}

func (r CatalogRepository) GetCatalog(ctx context.Context, storeID string) ([]*domain.CatalogProduct, error) {
	const query = "SELECT id, name, description, sku, price FROM %s WHERE store_id = $1"

	rows, err := r.db.QueryContext(ctx, r.table(query), storeID)
	if err != nil {
		return nil, fmt.Errorf("querying products: %w", err)
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			err = fmt.Errorf("closing product rows: %w", err)
		}
	}(rows)

	products := make([]*domain.CatalogProduct, 0)
	for rows.Next() {
		product := &domain.CatalogProduct{
			StoreID: storeID,
		}
		err := rows.Scan(&product.ID, &product.Name, &product.Description, &product.SKU, &product.Price)
		if err != nil {
			return nil, fmt.Errorf("scanning product: %w", err)
		}

		products = append(products, product)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("finishing product rows: %w", err)
	}

	return products, nil
}

func (r CatalogRepository) table(query string) string {
	return fmt.Sprintf(query, r.tableName)
}
