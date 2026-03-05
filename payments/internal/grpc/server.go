package grpc

import (
	"context"
	"mall/payments/internal/application"
	"mall/payments/paymentspb"

	"github.com/google/uuid"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc"
)

type server struct {
	app application.App
	paymentspb.UnimplementedPaymentsServiceServer
}

var _ paymentspb.PaymentsServiceServer = (*server)(nil)

func RegisterServer(app application.App, registrar grpc.ServiceRegistrar) error {
	paymentspb.RegisterPaymentsServiceServer(registrar, server{app: app})
	return nil
}

func (s server) AuthorizePayment(ctx context.Context, request *paymentspb.AuthorizePaymentRequest) (*paymentspb.AuthorizePaymentResponse, error) {
	span := trace.SpanFromContext(ctx)

	id := uuid.New().String()

	span.SetAttributes(
		attribute.String("PaymentID", id),
		attribute.String("CustomerID", request.GetCustomerId()),
		attribute.Float64("Amount", request.GetAmount()),
	)

	err := s.app.AuthorizePayment(ctx, application.AuthorizePayment{
		ID:         id,
		CustomerID: request.GetCustomerId(),
		Amount:     request.GetAmount(),
	})
	if err != nil {
		span.RecordError(err, trace.WithAttributes(attribute.String("error.message", err.Error())))
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	return &paymentspb.AuthorizePaymentResponse{Id: id}, nil
}

func (s server) ConfirmPayment(ctx context.Context, request *paymentspb.ConfirmPaymentRequest) (*paymentspb.ConfirmPaymentResponse, error) {
	span := trace.SpanFromContext(ctx)

	span.SetAttributes(
		attribute.String("PaymentID", request.GetId()),
	)

	err := s.app.ConfirmPayment(ctx, application.ConfirmPayment{
		ID: request.GetId(),
	})
	if err != nil {
		span.RecordError(err, trace.WithAttributes(attribute.String("error.message", err.Error())))
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	return &paymentspb.ConfirmPaymentResponse{}, nil
}

func (s server) CreateInvoice(ctx context.Context, request *paymentspb.CreateInvoiceRequest) (*paymentspb.CreateInvoiceResponse, error) {
	span := trace.SpanFromContext(ctx)

	id := uuid.New().String()

	span.SetAttributes(
		attribute.String("InvoiceID", id),
		attribute.String("OrderID", request.GetOrderId()),
		attribute.Float64("Amount", request.GetAmount()),
	)

	err := s.app.CreateInvoice(ctx, application.CreateInvoice{
		ID:      id,
		OrderID: request.GetOrderId(),
		Amount:  request.GetAmount(),
	})
	if err != nil {
		span.RecordError(err, trace.WithAttributes(attribute.String("error.message", err.Error())))
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	return &paymentspb.CreateInvoiceResponse{Id: id}, nil
}

func (s server) AdjustInvoice(ctx context.Context, request *paymentspb.AdjustInvoiceRequest) (*paymentspb.AdjustInvoiceResponse, error) {
	span := trace.SpanFromContext(ctx)

	span.SetAttributes(
		attribute.String("InvoiceID", request.GetId()),
		attribute.Float64("Amount", request.GetAmount()),
	)

	err := s.app.AdjustInvoice(ctx, application.AdjustInvoice{
		ID:     request.GetId(),
		Amount: request.GetAmount(),
	})
	if err != nil {
		span.RecordError(err, trace.WithAttributes(attribute.String("error.message", err.Error())))
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	return &paymentspb.AdjustInvoiceResponse{}, nil
}

func (s server) PayInvoice(ctx context.Context, request *paymentspb.PayInvoiceRequest) (*paymentspb.PayInvoiceResponse, error) {
	span := trace.SpanFromContext(ctx)

	span.SetAttributes(
		attribute.String("InvoiceID", request.GetId()),
	)

	err := s.app.PayInvoice(ctx, application.PayInvoice{
		ID: request.GetId(),
	})
	if err != nil {
		span.RecordError(err, trace.WithAttributes(attribute.String("error.message", err.Error())))
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	return &paymentspb.PayInvoiceResponse{}, nil
}

func (s server) CancelInvoice(ctx context.Context, request *paymentspb.CancelInvoiceRequest) (*paymentspb.CancelInvoiceResponse, error) {
	span := trace.SpanFromContext(ctx)

	span.SetAttributes(
		attribute.String("InvoiceID", request.GetId()),
	)

	err := s.app.CancelInvoice(ctx, application.CancelInvoice{
		ID: request.GetId(),
	})
	if err != nil {
		span.RecordError(err, trace.WithAttributes(attribute.String("error.message", err.Error())))
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	return &paymentspb.CancelInvoiceResponse{}, nil
}
