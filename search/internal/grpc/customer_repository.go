package grpc

import (
	"context"
	"mall/customers/customerspb"
	"mall/search/internal/domain"

	"google.golang.org/grpc"
)

type CustomerRepository struct {
	client customerspb.CustomersServiceClient
}

var _ domain.CustomerRepository = (*CustomerRepository)(nil)

func NewCustomerRepository(conn *grpc.ClientConn) CustomerRepository {
	return CustomerRepository{
		client: customerspb.NewCustomersServiceClient(conn),
	}
}

func (r CustomerRepository) Find(ctx context.Context, customerID string) (*domain.Customer, error) {
	resp, err := r.client.GetCustomer(ctx, &customerspb.GetCustomerRequest{
		Id: customerID,
	})
	if err != nil {
		return nil, err
	}

	return &domain.Customer{
		ID:   resp.GetCustomer().GetId(),
		Name: resp.GetCustomer().GetName(),
	}, nil
}
