package grpc

import (
	"context"
	"mall/depot/depotpb"
	"mall/depot/internal/application"
	"mall/depot/internal/application/commands"

	"github.com/google/uuid"
	"google.golang.org/grpc"
)

type server struct {
	app application.App
	depotpb.UnimplementedDepotServiceServer
}

var _ depotpb.DepotServiceServer = (*server)(nil)

func RegisterServer(app application.App, registrar grpc.ServiceRegistrar) error {
	depotpb.RegisterDepotServiceServer(registrar, server{app: app})
	return nil
}

func (s server) CreateShoppingList(ctx context.Context, request *depotpb.CreateShoppingListRequest) (*depotpb.CreateShoppingListResponse, error) {
	id := uuid.New().String()

	items := make([]commands.OrderItem, 0, len(request.GetItems()))
	for _, item := range request.GetItems() {
		items = append(items, s.itemToDomain(item))
	}

	err := s.app.CreateShoppingList(ctx, commands.CreateShoppingListRequest{
		ID:      id,
		OrderID: request.GetOrderId(),
		Items:   items,
	})
	if err != nil {
		return nil, err
	}

	return &depotpb.CreateShoppingListResponse{Id: id}, nil
}

func (s server) CancelShoppingList(ctx context.Context, request *depotpb.CancelShoppingListRequest) (*depotpb.CancelShoppingListResponse, error) {
	err := s.app.CancelShoppingList(ctx, commands.CancelShoppingListRequest{
		ID: request.GetId(),
	})
	if err != nil {
		return nil, err
	}

	return &depotpb.CancelShoppingListResponse{}, nil
}

func (s server) AssignShoppingList(ctx context.Context, request *depotpb.AssignShoppingListRequest) (*depotpb.AssignShoppingListResponse, error) {
	err := s.app.AssignShoppingList(ctx, commands.AssignShoppingListRequest{
		ID:    request.GetId(),
		BotID: request.GetBotId(),
	})
	if err != nil {
		return nil, err
	}

	return &depotpb.AssignShoppingListResponse{}, nil
}

func (s server) CompleteShoppingList(ctx context.Context, request *depotpb.CompleteShoppingListRequest) (*depotpb.CompleteShoppingListResponse, error) {
	err := s.app.CompleteShoppingList(ctx, commands.CompleteShoppingListRequest{
		ID: request.GetId(),
	})
	if err != nil {
		return nil, err
	}

	return &depotpb.CompleteShoppingListResponse{}, nil
}

func (s server) itemToDomain(item *depotpb.OrderItem) commands.OrderItem {
	return commands.OrderItem{
		StoreID:   item.GetStoreId(),
		ProductID: item.GetProductId(),
		Quantity:  int(item.GetQuantity()),
	}
}
