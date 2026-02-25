package grpc

import (
	"context"
	"fmt"
	"mall/baskets/internal/domain"
	"mall/internal/rpc"
	"mall/stores/storespb"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type ProductRepository struct {
	endpoint string
}

var _ domain.ProductRepository = (*ProductRepository)(nil)

func NewProductRepository(endpoint string) ProductRepository {
	return ProductRepository{
		endpoint: endpoint,
	}
}

func (r ProductRepository) Find(ctx context.Context, productID string) (*domain.Product, error) {
	conn, err := r.dial(ctx)
	if err != nil {
		return nil, err
	}

	defer func(conn *grpc.ClientConn) {
		conn.Close()
	}(conn)

	resp, err := storespb.NewStoresServiceClient(conn).GetProduct(ctx, &storespb.GetProductRequest{
		Id: productID,
	})
	if err != nil {
		if st, ok := status.FromError(err); ok && st.Code() == codes.NotFound {
			return nil, fmt.Errorf("could not locate product with given id")
		}

		return nil, fmt.Errorf("requesting product: %w", err)
	}

	return r.productToDomain(resp.Product), nil
}

func (r ProductRepository) productToDomain(product *storespb.Product) *domain.Product {
	return &domain.Product{
		ID:      product.GetId(),
		StoreID: product.GetStoreId(),
		Name:    product.GetName(),
		Price:   product.GetPrice(),
	}
}

func (r ProductRepository) dial(ctx context.Context) (*grpc.ClientConn, error) {
	return rpc.Dial(ctx, r.endpoint)
}
