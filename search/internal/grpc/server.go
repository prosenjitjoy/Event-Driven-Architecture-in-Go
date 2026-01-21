package grpc

import (
	"context"
	"mall/search/internal/application"
	"mall/search/searchpb"

	"google.golang.org/grpc"
)

type server struct {
	app application.Application
	searchpb.UnimplementedSearchServiceServer
}

func RegisterServer(ctx context.Context, app application.Application, registrar grpc.ServiceRegistrar) error {
	searchpb.RegisterSearchServiceServer(registrar, server{app: app})

	return nil
}

func (s server) SearchOrders(context.Context, *searchpb.SearchOrdersRequest) (*searchpb.SearchOrdersResponse, error) {
	panic("implement me")
}

func (s server) GetOrder(context.Context, *searchpb.GetOrderRequest) (*searchpb.GetOrderResponse, error) {
	panic("implement me")
}
