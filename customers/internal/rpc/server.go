package rpc

import (
	"context"

	"github.com/google/uuid"
	"google.golang.org/grpc"

	"mall/customers/customerspb"
	"mall/customers/internal/application"
	"mall/customers/internal/domain"
)

type server struct {
	app application.App
	customerspb.UnimplementedCustomersServiceServer
}

var _ customerspb.CustomersServiceServer = (*server)(nil)

func RegisterServer(app application.App, registrar grpc.ServiceRegistrar) error {
	customerspb.RegisterCustomersServiceServer(registrar, server{app: app})

	return nil
}

func (s server) RegisterCustomer(ctx context.Context, request *customerspb.RegisterCustomerRequest) (*customerspb.RegisterCustomerResponse, error) {
	id := uuid.NewString()

	err := s.app.RegisterCustomer(ctx, application.RegisterCustomer{
		ID:        id,
		Name:      request.GetName(),
		SmsNumber: request.GetSmsNumber(),
	})
	if err != nil {
		return nil, err
	}

	return &customerspb.RegisterCustomerResponse{Id: id}, nil
}

func (s server) AuthorizeCustomer(ctx context.Context, request *customerspb.AuthorizeCustomerRequest) (*customerspb.AuthorizeCustomerResponse, error) {
	err := s.app.AuthorizeCustomer(ctx, application.AuthorizeCustomer{
		ID: request.GetId(),
	})
	if err != nil {
		return nil, err
	}

	return &customerspb.AuthorizeCustomerResponse{}, nil
}

func (s server) GetCustomer(ctx context.Context, request *customerspb.GetCustomerRequest) (*customerspb.GetCustomerResponse, error) {
	customer, err := s.app.GetCustomer(ctx, application.GetCustomer{
		ID: request.GetId(),
	})
	if err != nil {
		return nil, err
	}

	return &customerspb.GetCustomerResponse{Customer: s.customerFromDomain(customer)}, nil
}

func (s server) EnableCustomer(ctx context.Context, request *customerspb.EnableCustomerRequest) (*customerspb.EnableCustomerResponse, error) {
	err := s.app.EnableCustomer(ctx, application.EnableCustomer{ID: request.GetId()})
	if err != nil {
		return nil, err
	}

	return &customerspb.EnableCustomerResponse{}, nil
}

func (s server) DisableCustomer(ctx context.Context, request *customerspb.DisableCustomerRequest) (*customerspb.DisableCustomerResponse, error) {
	err := s.app.DisableCustomer(ctx, application.DisableCustomer{ID: request.GetId()})
	if err != nil {
		return nil, err
	}

	return &customerspb.DisableCustomerResponse{}, nil
}

func (s server) customerFromDomain(customer *domain.Customer) *customerspb.Customer {
	return &customerspb.Customer{
		Id:        customer.ID,
		Name:      customer.Name,
		SmsNumber: customer.SmsNumber,
		Enabled:   customer.Enabled,
	}
}
