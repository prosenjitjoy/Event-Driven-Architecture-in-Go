package grpc

import (
	"context"
	"mall/customers/customerspb"
	"mall/internal/rpc"
	"mall/search/internal/domain"

	"google.golang.org/grpc"
)

type CustomerRepository struct {
	endpoint string
}

var _ domain.CustomerRepository = (*CustomerRepository)(nil)

func NewCustomerRepository(endpoint string) CustomerRepository {
	return CustomerRepository{
		endpoint: endpoint,
	}
}

func (r CustomerRepository) Find(ctx context.Context, customerID string) (*domain.Customer, error) {
	conn, err := r.dial(ctx)
	if err != nil {
		return nil, err
	}

	defer func(conn *grpc.ClientConn) {
		conn.Close()
	}(conn)

	resp, err := customerspb.NewCustomersServiceClient(conn).GetCustomer(ctx, &customerspb.GetCustomerRequest{Id: customerID})
	if err != nil {
		return nil, err
	}

	return &domain.Customer{
		ID:   resp.GetCustomer().GetId(),
		Name: resp.GetCustomer().GetName(),
	}, nil
}

func (r CustomerRepository) dial(ctx context.Context) (*grpc.ClientConn, error) {
	return rpc.Dial(ctx, r.endpoint)
}
