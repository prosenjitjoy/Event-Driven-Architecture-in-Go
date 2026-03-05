package grpc

import (
	"context"
	"mall/notifications/internal/application"
	"mall/notifications/notificationspb"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc"
)

type server struct {
	app application.App
	notificationspb.UnimplementedNotificationsServiceServer
}

var _ notificationspb.NotificationsServiceServer = (*server)(nil)

func RegisterServer(app application.App, registrar grpc.ServiceRegistrar) error {
	notificationspb.RegisterNotificationsServiceServer(registrar, server{app: app})

	return nil
}

func (s server) NotifyOrderCreated(ctx context.Context, request *notificationspb.NotifyOrderCreatedRequest) (*notificationspb.NotifyOrderCreatedResponse, error) {
	span := trace.SpanFromContext(ctx)

	span.SetAttributes(
		attribute.String("OrderID", request.GetOrderId()),
		attribute.String("CustomerID", request.GetCustomerId()),
	)

	err := s.app.NotifyOrderCreated(ctx, application.OrderCreated{
		OrderID:    request.GetOrderId(),
		CustomerID: request.GetCustomerId(),
	})
	if err != nil {
		span.RecordError(err, trace.WithAttributes(attribute.String("error.message", err.Error())))
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	return &notificationspb.NotifyOrderCreatedResponse{}, nil
}

func (s server) NotifyOrderCanceled(ctx context.Context, request *notificationspb.NotifyOrderCanceledRequest) (*notificationspb.NotifyOrderCanceledResponse, error) {
	span := trace.SpanFromContext(ctx)

	span.SetAttributes(
		attribute.String("OrderID", request.GetOrderId()),
		attribute.String("CustomerID", request.GetCustomerId()),
	)

	err := s.app.NotifyOrderCanceled(ctx, application.OrderCanceled{
		OrderID:    request.GetOrderId(),
		CustomerID: request.GetCustomerId(),
	})
	if err != nil {
		span.RecordError(err, trace.WithAttributes(attribute.String("error.message", err.Error())))
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	return &notificationspb.NotifyOrderCanceledResponse{}, nil
}

func (s server) NotifyOrderReady(ctx context.Context, request *notificationspb.NotifyOrderReadyRequest) (*notificationspb.NotifyOrderReadyResponse, error) {
	span := trace.SpanFromContext(ctx)

	span.SetAttributes(
		attribute.String("OrderID", request.GetOrderId()),
		attribute.String("CustomerID", request.GetCustomerId()),
	)

	err := s.app.NotifyOrderReady(ctx, application.OrderReady{
		OrderID:    request.GetOrderId(),
		CustomerID: request.GetCustomerId(),
	})
	if err != nil {
		span.RecordError(err, trace.WithAttributes(attribute.String("error.message", err.Error())))
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	return &notificationspb.NotifyOrderReadyResponse{}, nil
}
