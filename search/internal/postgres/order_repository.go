package postgres

import (
	"bytes"
	"context"
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"mall/internal/postgres"
	"mall/search/internal/domain"
	"strings"
)

type OrderRepository struct {
	tableName string
	db        postgres.DBTX
}

var _ domain.OrderRepository = (*OrderRepository)(nil)

func NewOrderRepository(tableName string, db postgres.DBTX) OrderRepository {
	return OrderRepository{
		tableName: tableName,
		db:        db,
	}
}

func (r OrderRepository) Add(ctx context.Context, order *domain.Order) error {
	const query = "INSERT INTO %s (order_id, customer_id, customer_name, items, status, product_ids, store_ids, created_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)"

	items, err := json.Marshal(order.Items)
	if err != nil {
		return err
	}

	productIDs := make(IDArray, len(order.Items))
	storeMap := make(map[string]struct{})

	for i, item := range order.Items {
		productIDs[i] = item.ProductID
		storeMap[item.StoreID] = struct{}{}
	}

	storeIDs := make(IDArray, 0, len(storeMap))
	for storeID, _ := range storeMap {
		storeIDs = append(storeIDs, storeID)
	}

	_, err = r.db.ExecContext(ctx, r.table(query), order.OrderID, order.CustomerID, order.CustomerName, items, order.Status, productIDs, storeIDs, order.CreatedAt)
	if err != nil {
		return err
	}

	return nil
}

func (r OrderRepository) UpdateStatus(ctx context.Context, orderID, status string) error {
	const query = "UPDATE %s SET status = $2 WHERE id = $1"

	_, err := r.db.ExecContext(ctx, r.table(query), orderID, status)
	if err != nil {
		return err
	}

	return nil
}

func (r OrderRepository) Search(ctx context.Context, search domain.SearchOrders) ([]domain.Order, error) {
	panic("implement me")
}

func (r OrderRepository) Get(ctx context.Context, orderID string) (*domain.Order, error) {
	const query = "SELECT customer_id, customer_name, items, status, created_at FROM %s WHERE order_id = $1"

	order := &domain.Order{
		OrderID: orderID,
	}

	var itemData []byte
	err := r.db.QueryRowContext(ctx, r.table(query), orderID).Scan(&order.CustomerID, &order.CustomerName, &itemData, &order.Status, &order.CreatedAt)
	if err != nil {
		return nil, err
	}

	var items []domain.Item
	err = json.Unmarshal(itemData, &items)
	if err != nil {
		return nil, err
	}

	order.Items = items

	return order, nil
}

func (r OrderRepository) table(query string) string {
	return fmt.Sprintf(query, r.tableName)
}

type IDArray []string

func (a *IDArray) Scan(src any) error {
	var sep = []byte(",")

	var data []byte
	switch v := src.(type) {
	case string:
		data = []byte(v)
	case []byte:
		data = v
	default:
		return fmt.Errorf("INVALID_ARGUMENT. IDArray: unsupported type: %T", src)
	}

	ids := make([]string, bytes.Count(data, sep))
	for i, id := range bytes.Split(bytes.Trim(data, "{}"), sep) {
		ids[i] = string(id)
	}

	*a = ids

	return nil
}

func (a IDArray) Value() (driver.Value, error) {
	if a == nil {
		return nil, nil
	}
	if len(a) == 0 {
		return "{}", nil
	}

	return fmt.Sprintf("{%s}", strings.Join(a, ",")), nil
}
