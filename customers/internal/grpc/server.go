package grpc

import (
	"context"
	"mall/customers/customerspb"
	"mall/customers/internal/application"
	"mall/customers/internal/domain"

	"github.com/google/uuid"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc"
)

type server struct {
	app application.App
	customerspb.UnimplementedCustomersServiceServer
}

var _ customerspb.CustomersServiceServer = (*server)(nil)

func RegisterServer(app application.App, registrar grpc.ServiceRegistrar) error {
	customerspb.RegisterCustomersServiceServer(registrar, server{
		app: app,
	})

	return nil
}

func (s server) RegisterCustomer(ctx context.Context, request *customerspb.RegisterCustomerRequest) (*customerspb.RegisterCustomerResponse, error) {
	span := trace.SpanFromContext(ctx)

	id := uuid.New().String()

	span.SetAttributes(
		attribute.String("CustomerID", id),
		attribute.String("Name", request.GetName()),
		attribute.String("SmsNumber", request.GetSmsNumber()),
	)

	err := s.app.RegisterCustomer(ctx, application.RegisterCustomer{
		ID:        id,
		Name:      request.GetName(),
		SmsNumber: request.GetSmsNumber(),
	})
	if err != nil {
		span.RecordError(err, trace.WithAttributes(attribute.String("error.message", err.Error())))
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	return &customerspb.RegisterCustomerResponse{Id: id}, nil
}

func (s server) AuthorizeCustomer(ctx context.Context, request *customerspb.AuthorizeCustomerRequest) (*customerspb.AuthorizeCustomerResponse, error) {
	span := trace.SpanFromContext(ctx)

	span.SetAttributes(
		attribute.String("CustomerID", request.GetId()),
	)

	err := s.app.AuthorizeCustomer(ctx, application.AuthorizeCustomer{
		ID: request.GetId(),
	})
	if err != nil {
		span.RecordError(err, trace.WithAttributes(attribute.String("error.message", err.Error())))
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	return &customerspb.AuthorizeCustomerResponse{}, nil
}

func (s server) GetCustomer(ctx context.Context, request *customerspb.GetCustomerRequest) (*customerspb.GetCustomerResponse, error) {
	span := trace.SpanFromContext(ctx)

	span.SetAttributes(
		attribute.String("CustomerID", request.GetId()),
	)

	customer, err := s.app.GetCustomer(ctx, application.GetCustomer{
		ID: request.GetId(),
	})
	if err != nil {
		span.RecordError(err, trace.WithAttributes(attribute.String("error.message", err.Error())))
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	return &customerspb.GetCustomerResponse{
		Customer: s.customerFromDomain(customer),
	}, nil
}

func (s server) EnableCustomer(ctx context.Context, request *customerspb.EnableCustomerRequest) (*customerspb.EnableCustomerResponse, error) {
	span := trace.SpanFromContext(ctx)

	span.SetAttributes(
		attribute.String("CustomerID", request.GetId()),
	)

	err := s.app.EnableCustomer(ctx, application.EnableCustomer{
		ID: request.GetId(),
	})
	if err != nil {
		span.RecordError(err, trace.WithAttributes(attribute.String("error.message", err.Error())))
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	return &customerspb.EnableCustomerResponse{}, nil
}

func (s server) DisableCustomer(ctx context.Context, request *customerspb.DisableCustomerRequest) (*customerspb.DisableCustomerResponse, error) {
	span := trace.SpanFromContext(ctx)

	span.SetAttributes(
		attribute.String("CustomerID", request.GetId()),
	)

	err := s.app.DisableCustomer(ctx, application.DisableCustomer{
		ID: request.GetId(),
	})
	if err != nil {
		span.RecordError(err, trace.WithAttributes(attribute.String("error.message", err.Error())))
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	return &customerspb.DisableCustomerResponse{}, nil
}

func (s server) customerFromDomain(customer *domain.Customer) *customerspb.Customer {
	return &customerspb.Customer{
		Id:        customer.ID(),
		Name:      customer.Name,
		SmsNumber: customer.SmsNumber,
		Enabled:   customer.Enabled,
	}
}
