//go:build integration

package grpc

import (
	"context"
	"mall/baskets/basketspb"
	"mall/baskets/internal/application"
	"mall/baskets/internal/domain"
	"mall/internal/ddd"
	"mall/internal/es"
	"net"
	"testing"

	"github.com/stretchr/testify/suite"
	"go.uber.org/mock/gomock"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type MockApplication struct {
	baskets   *domain.MockBasketRepository
	stores    *domain.MockStoreRepository
	products  *domain.MockProductRepository
	publisher *ddd.MockEventPublisher[ddd.Event]
}

type serverSuite struct {
	mockApp *MockApplication
	server  *grpc.Server
	client  basketspb.BasketServiceClient
	ctrl    *gomock.Controller
	suite.Suite
}

func TestServer(t *testing.T) {
	suite.Run(t, &serverSuite{})
}

func (s *serverSuite) SetupSuite()    { s.ctrl = gomock.NewController(s.T()) }
func (s *serverSuite) TearDownSuite() { s.ctrl.Finish() }

func (s *serverSuite) SetupTest() {
	const grpcTestPort = ":10912"

	s.server = grpc.NewServer()

	listener, err := net.Listen("tcp", grpcTestPort)
	s.Require().NoError(err)

	s.mockApp = &MockApplication{
		baskets:   domain.NewMockBasketRepository(s.ctrl),
		stores:    domain.NewMockStoreRepository(s.ctrl),
		products:  domain.NewMockProductRepository(s.ctrl),
		publisher: ddd.NewMockEventPublisher[ddd.Event](s.ctrl),
	}

	app := application.New(s.mockApp.baskets, s.mockApp.stores, s.mockApp.products, s.mockApp.publisher)

	err = RegisterServer(app, s.server)
	s.Require().NoError(err)

	go func(listener net.Listener) {
		err := s.server.Serve(listener)
		s.Require().NoError(err)
	}(listener)

	// create client
	conn, err := grpc.NewClient(grpcTestPort, grpc.WithTransportCredentials(insecure.NewCredentials()))
	s.Require().NoError(err)

	s.client = basketspb.NewBasketServiceClient(conn)
}
func (s *serverSuite) TearDownTest() {
	s.server.GracefulStop()
}

func (s *serverSuite) TestBasketService_StartBasket() {
	basket := &domain.Basket{
		Aggregate: es.NewAggregate("basket-id", domain.BasketAggregate),
	}

	s.mockApp.baskets.EXPECT().Load(gomock.Any(), gomock.Any()).Times(1).Return(basket, nil)

	s.mockApp.baskets.EXPECT().Save(gomock.Any(), gomock.AssignableToTypeOf(&domain.Basket{})).Times(1).Return(nil)

	s.mockApp.publisher.EXPECT().Publish(gomock.Any(), gomock.Any()).Times(1).Return(nil)

	_, err := s.client.StartBasket(context.Background(), &basketspb.StartBasketRequest{CustomerId: "customer-id"})
	s.Assert().NoError(err)
}

func (s *serverSuite) TestBasketService_CancelBasket() {
	basket := &domain.Basket{
		Aggregate:  es.NewAggregate("basket-id", domain.BasketAggregate),
		CustomerID: "customer-id",
		Items:      make(map[string]domain.Item),
		Status:     domain.BasketIsOpen,
	}

	s.mockApp.baskets.EXPECT().Load(gomock.Any(), gomock.Eq("basket-id")).Times(1).Return(basket, nil)

	s.mockApp.baskets.EXPECT().Save(gomock.Any(), gomock.AssignableToTypeOf(&domain.Basket{})).Times(1).Return(nil)

	s.mockApp.publisher.EXPECT().Publish(gomock.Any(), gomock.Any()).Times(1).Return(nil)

	_, err := s.client.CancelBasket(context.Background(), &basketspb.CancelBasketRequest{Id: "basket-id"})
	s.Assert().NoError(err)
}

func (s *serverSuite) TestBasketService_CheckoutBasket() {
	basket := &domain.Basket{
		Aggregate:  es.NewAggregate("basket-id", domain.BasketAggregate),
		CustomerID: "customer-id",
		PaymentID:  "payment-id",
		Items: map[string]domain.Item{
			"product-id": {
				StoreID:      "store-id",
				ProductID:    "product-id",
				StoreName:    "store-name",
				ProductName:  "product-name",
				ProductPrice: 1.00,
				Quantity:     1,
			},
		},
		Status: domain.BasketIsOpen,
	}

	s.mockApp.baskets.EXPECT().Load(gomock.Any(), gomock.Eq("basket-id")).Times(1).Return(basket, nil)

	s.mockApp.baskets.EXPECT().Save(gomock.Any(), gomock.AssignableToTypeOf(&domain.Basket{})).Times(1).Return(nil)

	s.mockApp.publisher.EXPECT().Publish(gomock.Any(), gomock.Any()).Times(1).Return(nil)

	_, err := s.client.CheckoutBasket(context.Background(), &basketspb.CheckoutBasketRequest{
		Id:        "basket-id",
		PaymentId: "payment-id",
	})
	s.Assert().NoError(err)
}

func (s *serverSuite) TestBasketService_AddItem() {
	store := &domain.Store{
		ID:   "store-id",
		Name: "store-name",
	}

	product := &domain.Product{
		ID:      "product-id",
		StoreID: "store-id",
		Name:    "product-name",
		Price:   10.00,
	}

	basket := &domain.Basket{
		Aggregate:  es.NewAggregate("basket-id", domain.BasketAggregate),
		CustomerID: "customer-id",
		PaymentID:  "payment-id",
		Items:      make(map[string]domain.Item),
		Status:     domain.BasketIsOpen,
	}

	s.mockApp.baskets.EXPECT().Load(gomock.Any(), gomock.Eq("basket-id")).Times(1).Return(basket, nil)

	s.mockApp.products.EXPECT().Find(gomock.Any(), gomock.Eq("product-id")).Times(1).Return(product, nil)

	s.mockApp.stores.EXPECT().Find(gomock.Any(), gomock.Eq("store-id")).Times(1).Return(store, nil)

	s.mockApp.baskets.EXPECT().Save(gomock.Any(), gomock.AssignableToTypeOf(&domain.Basket{})).Times(1).Return(nil)

	_, err := s.client.AddItem(context.Background(), &basketspb.AddItemRequest{
		Id:        "basket-id",
		ProductId: "product-id",
		Quantity:  1,
	})
	s.Assert().NoError(err)
}

func (s *serverSuite) TestBasketService_RemoveItem() {
	store := &domain.Store{
		ID:   "store-id",
		Name: "store-name",
	}

	product := &domain.Product{
		ID:      "product-id",
		StoreID: "store-id",
		Name:    "product-name",
		Price:   10.00,
	}

	item := domain.Item{
		StoreID:      store.ID,
		ProductID:    product.ID,
		StoreName:    store.Name,
		ProductName:  product.Name,
		ProductPrice: product.Price,
		Quantity:     10,
	}

	s.mockApp.products.EXPECT().Find(gomock.Any(), gomock.Eq("product-id")).Times(1).Return(product, nil)

	basket := &domain.Basket{
		Aggregate:  es.NewAggregate("basket-id", domain.BasketAggregate),
		CustomerID: "customer-id",
		PaymentID:  "payment-id",
		Items:      map[string]domain.Item{product.ID: item},
		Status:     domain.BasketIsOpen,
	}

	s.mockApp.baskets.EXPECT().Load(gomock.Any(), gomock.Eq("basket-id")).Times(1).Return(basket, nil)

	s.mockApp.baskets.EXPECT().Save(gomock.Any(), gomock.AssignableToTypeOf(&domain.Basket{})).Times(1).Return(nil)

	_, err := s.client.RemoveItem(context.Background(), &basketspb.RemoveItemRequest{
		Id:        "basket-id",
		ProductId: "product-id",
		Quantity:  1,
	})
	s.Assert().NoError(err)
}

func (s *serverSuite) TestBasketService_GetBasket() {
	basket := &domain.Basket{
		Aggregate:  es.NewAggregate("basket-id", domain.BasketAggregate),
		CustomerID: "customer-id",
		PaymentID:  "payment-id",
		Items: map[string]domain.Item{
			"product-id": {
				StoreID:      "store-id",
				ProductID:    "product-id",
				StoreName:    "store-name",
				ProductName:  "product-name",
				ProductPrice: 1.00,
				Quantity:     1,
			},
		},
		Status: domain.BasketIsOpen,
	}

	s.mockApp.baskets.EXPECT().Load(gomock.Any(), gomock.Eq("basket-id")).Times(1).Return(basket, nil)

	resp, err := s.client.GetBasket(context.Background(), &basketspb.GetBasketRequest{Id: "basket-id"})
	s.Require().NoError(err)
	s.Assert().Equal(basket.ID(), resp.Basket.GetId())
}
