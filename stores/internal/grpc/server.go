package grpc

import (
	"context"
	"mall/stores/internal/application"
	"mall/stores/internal/application/commands"
	"mall/stores/internal/application/queries"
	"mall/stores/internal/domain"
	"mall/stores/storespb"

	"github.com/google/uuid"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc"
)

type server struct {
	app application.App
	storespb.UnimplementedStoresServiceServer
}

var _ storespb.StoresServiceServer = (*server)(nil)

func RegisterServer(app application.App, registrar grpc.ServiceRegistrar) error {
	storespb.RegisterStoresServiceServer(registrar, server{app: app})
	return nil
}

func (s server) CreateStore(ctx context.Context, request *storespb.CreateStoreRequest) (*storespb.CreateStoreResponse, error) {
	span := trace.SpanFromContext(ctx)

	storeID := uuid.New().String()

	span.SetAttributes(
		attribute.String("StoreID", storeID),
		attribute.String("Name", request.GetName()),
		attribute.String("Location", request.GetLocation()),
	)

	err := s.app.CreateStore(ctx, commands.CreateStoreRequest{
		ID:       storeID,
		Name:     request.GetName(),
		Location: request.GetLocation(),
	})
	if err != nil {
		span.RecordError(err, trace.WithAttributes(attribute.String("error.message", err.Error())))
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	return &storespb.CreateStoreResponse{
		Id: storeID,
	}, nil
}

func (s server) EnableParticipation(ctx context.Context, request *storespb.EnableParticipationRequest) (*storespb.EnableParticipationResponse, error) {
	span := trace.SpanFromContext(ctx)

	span.SetAttributes(
		attribute.String("StoreID", request.GetId()),
	)

	err := s.app.EnableParticipation(ctx, commands.EnableParticipationRequest{
		ID: request.GetId(),
	})
	if err != nil {
		span.RecordError(err, trace.WithAttributes(attribute.String("error.message", err.Error())))
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	return &storespb.EnableParticipationResponse{}, nil
}

func (s server) DisableParticipation(ctx context.Context, request *storespb.DisableParticipationRequest) (*storespb.DisableParticipationResponse, error) {
	span := trace.SpanFromContext(ctx)

	span.SetAttributes(
		attribute.String("StoreID", request.GetId()),
	)

	err := s.app.DisableParticipation(ctx, commands.DisableParticipationRequest{
		ID: request.GetId(),
	})
	if err != nil {
		span.RecordError(err, trace.WithAttributes(attribute.String("error.message", err.Error())))
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	return &storespb.DisableParticipationResponse{}, nil
}

func (s server) RebrandStore(ctx context.Context, request *storespb.RebrandStoreRequest) (*storespb.RebrandStoreResponse, error) {
	span := trace.SpanFromContext(ctx)

	span.SetAttributes(
		attribute.String("StoreID", request.GetId()),
		attribute.String("NewName", request.GetName()),
	)

	err := s.app.RebrandStore(ctx, commands.RebrandStoreRequest{
		ID:   request.GetId(),
		Name: request.GetName(),
	})
	if err != nil {
		span.RecordError(err, trace.WithAttributes(attribute.String("error.message", err.Error())))
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	return &storespb.RebrandStoreResponse{}, nil
}

func (s server) GetStore(ctx context.Context, request *storespb.GetStoreRequest) (*storespb.GetStoreResponse, error) {
	span := trace.SpanFromContext(ctx)

	span.SetAttributes(
		attribute.String("StoreID", request.GetId()),
	)

	store, err := s.app.GetStore(ctx, queries.GetStoreRequest{
		ID: request.GetId(),
	})
	if err != nil {
		span.RecordError(err, trace.WithAttributes(attribute.String("error.message", err.Error())))
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	return &storespb.GetStoreResponse{
		Store: s.storeFromDomain(store),
	}, nil
}

func (s server) GetStores(ctx context.Context, request *storespb.GetStoresRequest) (*storespb.GetStoresResponse, error) {
	span := trace.SpanFromContext(ctx)

	stores, err := s.app.GetStores(ctx, queries.GetStoresRequest{})
	if err != nil {
		span.RecordError(err, trace.WithAttributes(attribute.String("error.message", err.Error())))
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	protoStores := []*storespb.Store{}
	for _, store := range stores {
		protoStores = append(protoStores, s.storeFromDomain(store))
	}

	return &storespb.GetStoresResponse{
		Stores: protoStores,
	}, nil
}

func (s server) GetParticipatingStores(ctx context.Context, request *storespb.GetParticipatingStoresRequest) (*storespb.GetParticipatingStoresResponse, error) {
	span := trace.SpanFromContext(ctx)

	stores, err := s.app.GetParticipatingStores(ctx, queries.GetParticipatingStoreRequest{})
	if err != nil {
		span.RecordError(err, trace.WithAttributes(attribute.String("error.message", err.Error())))
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	protoStores := []*storespb.Store{}
	for _, store := range stores {
		protoStores = append(protoStores, s.storeFromDomain(store))
	}

	return &storespb.GetParticipatingStoresResponse{
		Stores: protoStores,
	}, nil
}

func (s server) AddProduct(ctx context.Context, request *storespb.AddProductRequest) (*storespb.AddProductResponse, error) {
	span := trace.SpanFromContext(ctx)

	id := uuid.New().String()

	span.SetAttributes(
		attribute.String("ProductID", id),
	)

	err := s.app.AddProduct(ctx, commands.AddProductRequest{
		ID:          id,
		StoreID:     request.GetStoreId(),
		Name:        request.GetName(),
		Description: request.GetDescription(),
		SKU:         request.GetSku(),
		Price:       request.GetPrice(),
	})
	if err != nil {
		span.RecordError(err, trace.WithAttributes(attribute.String("error.message", err.Error())))
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	return &storespb.AddProductResponse{Id: id}, nil
}

func (s server) RebrandProduct(ctx context.Context, request *storespb.RebrandProductRequest) (*storespb.RebrandProductResponse, error) {
	span := trace.SpanFromContext(ctx)

	span.SetAttributes(
		attribute.String("ProductID", request.GetId()),
		attribute.String("Name", request.GetName()),
	)

	err := s.app.RebrandProduct(ctx, commands.RebrandProductRequest{
		ID:          request.GetId(),
		Name:        request.GetName(),
		Description: request.GetDescription(),
	})
	if err != nil {
		span.RecordError(err, trace.WithAttributes(attribute.String("error.message", err.Error())))
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	return &storespb.RebrandProductResponse{}, nil
}

func (s server) IncreaseProductPrice(ctx context.Context, request *storespb.IncreaseProductPriceRequest) (*storespb.IncreaseProductPriceResponse, error) {
	span := trace.SpanFromContext(ctx)

	span.SetAttributes(
		attribute.String("ProductID", request.GetId()),
	)

	err := s.app.IncreaseProductPrice(ctx, commands.IncreaseProductPriceRequest{
		ID:    request.GetId(),
		Price: request.GetPrice(),
	})
	if err != nil {
		span.RecordError(err, trace.WithAttributes(attribute.String("error.message", err.Error())))
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	return &storespb.IncreaseProductPriceResponse{}, nil
}

func (s server) DecreaseProductPrice(ctx context.Context, request *storespb.DecreaseProductPriceRequest) (*storespb.DecreaseProductPriceResponse, error) {
	span := trace.SpanFromContext(ctx)

	span.SetAttributes(
		attribute.String("ProductID", request.GetId()),
	)

	err := s.app.DecreaseProductPrice(ctx, commands.DecreaseProductPriceRequest{
		ID:    request.GetId(),
		Price: request.GetPrice(),
	})
	if err != nil {
		span.RecordError(err, trace.WithAttributes(attribute.String("error.message", err.Error())))
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	return &storespb.DecreaseProductPriceResponse{}, nil
}

func (s server) RemoveProduct(ctx context.Context, request *storespb.RemoveProductRequest) (*storespb.RemoveProductResponse, error) {
	span := trace.SpanFromContext(ctx)

	span.SetAttributes(
		attribute.String("ProductID", request.GetId()),
	)

	err := s.app.RemoveProduct(ctx, commands.RemoveProductRequest{
		ID: request.GetId(),
	})
	if err != nil {
		span.RecordError(err, trace.WithAttributes(attribute.String("error.message", err.Error())))
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	return &storespb.RemoveProductResponse{}, nil
}

func (s server) GetProduct(ctx context.Context, request *storespb.GetProductRequest) (*storespb.GetProductResponse, error) {
	span := trace.SpanFromContext(ctx)

	span.SetAttributes(
		attribute.String("ProductID", request.GetId()),
	)

	product, err := s.app.GetProduct(ctx, queries.GetProductRequest{
		ID: request.GetId(),
	})
	if err != nil {
		span.RecordError(err, trace.WithAttributes(attribute.String("error.message", err.Error())))
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	return &storespb.GetProductResponse{
		Product: s.productFromDomain(product),
	}, nil
}

func (s server) GetCatalog(ctx context.Context, request *storespb.GetCatalogRequest) (*storespb.GetCatalogResponse, error) {
	span := trace.SpanFromContext(ctx)

	span.SetAttributes(
		attribute.String("StoreID", request.GetStoreId()),
	)

	products, err := s.app.GetCatalog(ctx, queries.GetCatalogRequest{
		StoreID: request.GetStoreId(),
	})
	if err != nil {
		span.RecordError(err, trace.WithAttributes(attribute.String("error.message", err.Error())))
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	protoProducts := make([]*storespb.Product, len(products))
	for i, product := range products {
		protoProducts[i] = s.productFromDomain(product)
	}

	return &storespb.GetCatalogResponse{
		Products: protoProducts,
	}, nil
}

func (s server) storeFromDomain(store *domain.MallStore) *storespb.Store {
	return &storespb.Store{
		Id:            store.ID,
		Name:          store.Name,
		Location:      store.Location,
		Participating: store.Participating,
	}
}

func (s server) productFromDomain(product *domain.CatalogProduct) *storespb.Product {
	return &storespb.Product{
		Id:          product.ID,
		StoreId:     product.StoreID,
		Name:        product.Name,
		Description: product.Description,
		Sku:         product.SKU,
		Price:       product.Price,
	}
}
