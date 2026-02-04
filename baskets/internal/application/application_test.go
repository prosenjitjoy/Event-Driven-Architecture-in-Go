package application

import (
	"context"
	"fmt"
	"mall/baskets/internal/domain"
	"mall/internal/ddd"
	"mall/internal/es"

	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestApplication_AddItem(t *testing.T) {
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

	type MockApplication struct {
		baskets   *domain.MockBasketRepository
		stores    *domain.MockStoreRepository
		products  *domain.MockProductRepository
		publisher *ddd.MockEventPublisher[ddd.Event]
	}

	type args struct {
		ctx context.Context
		add AddItem
	}

	testCases := []struct {
		name       string
		args       args
		buildStubs func(app MockApplication)
		wantErr    bool
	}{
		{
			name: "OK",
			args: args{
				ctx: context.Background(),
				add: AddItem{
					ID:        "basket-id",
					ProductID: "product-id",
					Quantity:  1,
				},
			},
			buildStubs: func(app MockApplication) {
				basket := &domain.Basket{
					Aggregate:  es.NewAggregate("basket-id", domain.BasketAggregate),
					CustomerID: "customer-id",
					PaymentID:  "payment-id",
					Items:      make(map[string]domain.Item),
					Status:     domain.BasketIsOpen,
				}

				app.baskets.EXPECT().Load(gomock.Any(), gomock.Eq("basket-id")).Times(1).Return(basket, nil)

				app.products.EXPECT().Find(gomock.Any(), gomock.Eq("product-id")).Times(1).Return(product, nil)

				app.stores.EXPECT().Find(gomock.Any(), gomock.Eq("store-id")).Times(1).Return(store, nil)

				app.baskets.EXPECT().Save(gomock.Any(), gomock.AssignableToTypeOf(&domain.Basket{})).Times(1).Return(nil)
			},
			wantErr: false,
		},
		{
			name: "NoBasket",
			args: args{
				ctx: context.Background(),
				add: AddItem{
					ID:        "basket-id",
					ProductID: "product-id",
					Quantity:  1,
				},
			},
			buildStubs: func(app MockApplication) {
				app.baskets.EXPECT().Load(gomock.Any(), gomock.Eq("basket-id")).Times(1).Return(nil, fmt.Errorf("no basket"))

				app.products.EXPECT().Find(gomock.Any(), gomock.Any()).Times(0)

				app.stores.EXPECT().Find(gomock.Any(), gomock.Any()).Times(0)

				app.baskets.EXPECT().Save(gomock.Any(), gomock.Any()).Times(0)

			},
			wantErr: true,
		},
		{
			name: "NoProduct",
			args: args{
				ctx: context.Background(),
				add: AddItem{
					ID:        "basket-id",
					ProductID: "product-id",
					Quantity:  1,
				},
			},
			buildStubs: func(app MockApplication) {
				basket := &domain.Basket{
					Aggregate:  es.NewAggregate("basket-id", domain.BasketAggregate),
					CustomerID: "customer-id",
					PaymentID:  "payment-id",
					Items:      make(map[string]domain.Item),
					Status:     domain.BasketIsOpen,
				}

				app.baskets.EXPECT().Load(gomock.Any(), gomock.Eq("basket-id")).Times(1).Return(basket, nil)

				app.products.EXPECT().Find(gomock.Any(), gomock.Eq("product-id")).Times(1).Return(nil, fmt.Errorf("no product"))

				app.stores.EXPECT().Find(gomock.Any(), gomock.Any()).Times(0)

				app.baskets.EXPECT().Save(gomock.Any(), gomock.Any()).Times(0)

			},
			wantErr: true,
		},
		{
			name: "NoStore",
			args: args{
				ctx: context.Background(),
				add: AddItem{
					ID:        "basket-id",
					ProductID: "product-id",
					Quantity:  1,
				},
			},
			buildStubs: func(app MockApplication) {
				basket := &domain.Basket{
					Aggregate:  es.NewAggregate("basket-id", domain.BasketAggregate),
					CustomerID: "customer-id",
					PaymentID:  "payment-id",
					Items:      make(map[string]domain.Item),
					Status:     domain.BasketIsOpen,
				}

				app.baskets.EXPECT().Load(gomock.Any(), gomock.Eq("basket-id")).Times(1).Return(basket, nil)

				app.products.EXPECT().Find(gomock.Any(), gomock.Eq("product-id")).Times(1).Return(product, nil)

				app.stores.EXPECT().Find(gomock.Any(), gomock.Eq("store-id")).Times(1).Return(nil, fmt.Errorf("no store"))

				app.baskets.EXPECT().Save(gomock.Any(), gomock.Any()).Times(0)

			},
			wantErr: true,
		},
		{
			name: "SaveFailed",
			args: args{
				ctx: context.Background(),
				add: AddItem{
					ID:        "basket-id",
					ProductID: "product-id",
					Quantity:  1,
				},
			},
			buildStubs: func(app MockApplication) {
				basket := &domain.Basket{
					Aggregate:  es.NewAggregate("basket-id", domain.BasketAggregate),
					CustomerID: "customer-id",
					PaymentID:  "payment-id",
					Items:      make(map[string]domain.Item),
					Status:     domain.BasketIsOpen,
				}

				app.baskets.EXPECT().Load(gomock.Any(), gomock.Eq("basket-id")).Times(1).Return(basket, nil)

				app.products.EXPECT().Find(gomock.Any(), gomock.Eq("product-id")).Times(1).Return(product, nil)

				app.stores.EXPECT().Find(gomock.Any(), gomock.Eq("store-id")).Times(1).Return(store, nil)

				app.baskets.EXPECT().Save(gomock.Any(), gomock.AssignableToTypeOf(&domain.Basket{})).Times(1).Return(fmt.Errorf("save failed"))

			},
			wantErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			m := MockApplication{
				baskets:   domain.NewMockBasketRepository(ctrl),
				stores:    domain.NewMockStoreRepository(ctrl),
				products:  domain.NewMockProductRepository(ctrl),
				publisher: ddd.NewMockEventPublisher[ddd.Event](ctrl),
			}

			app := New(m.baskets, m.stores, m.products, m.publisher)

			tc.buildStubs(m)

			err := app.AddItem(tc.args.ctx, tc.args.add)
			assert.Equal(t, err != nil, tc.wantErr)
		})
	}
}

func TestApplication_CancelBasket(t *testing.T) {
	type MockApplication struct {
		baskets   *domain.MockBasketRepository
		stores    *domain.MockStoreRepository
		products  *domain.MockProductRepository
		publisher *ddd.MockEventPublisher[ddd.Event]
	}

	type args struct {
		ctx    context.Context
		cancel CancelBasket
	}

	testCases := []struct {
		name       string
		args       args
		buildStubs func(app MockApplication)
		wantErr    bool
	}{
		{
			name: "OK",
			args: args{
				ctx:    context.Background(),
				cancel: CancelBasket{ID: "basket-id"},
			},
			buildStubs: func(app MockApplication) {
				basket := &domain.Basket{
					Aggregate:  es.NewAggregate("basket-id", domain.BasketAggregate),
					CustomerID: "customer-id",
					PaymentID:  "payment-id",
					Items:      make(map[string]domain.Item),
					Status:     domain.BasketIsOpen,
				}

				app.baskets.EXPECT().Load(gomock.Any(), gomock.Eq("basket-id")).Times(1).Return(basket, nil)

				app.baskets.EXPECT().Save(gomock.Any(), gomock.AssignableToTypeOf(&domain.Basket{})).Times(1).Return(nil)

				app.publisher.EXPECT().Publish(gomock.Any(), gomock.Any()).Times(1).Return(nil)

			},
			wantErr: false,
		},
		{
			name: "NoBasket",
			args: args{
				ctx:    context.Background(),
				cancel: CancelBasket{ID: "basket-id"},
			},
			buildStubs: func(app MockApplication) {
				app.baskets.EXPECT().Load(gomock.Any(), gomock.Eq("basket-id")).Times(1).Return(nil, fmt.Errorf("no basket"))

				app.baskets.EXPECT().Save(gomock.Any(), gomock.Any()).Times(0)

				app.publisher.EXPECT().Publish(gomock.Any(), gomock.Any()).Times(0)
			},
			wantErr: true,
		},
		{
			name: "SaveFailed",
			args: args{
				ctx:    context.Background(),
				cancel: CancelBasket{ID: "basket-id"},
			},
			buildStubs: func(app MockApplication) {
				basket := &domain.Basket{
					Aggregate:  es.NewAggregate("basket-id", domain.BasketAggregate),
					CustomerID: "customer-id",
					PaymentID:  "payment-id",
					Items:      make(map[string]domain.Item),
					Status:     domain.BasketIsOpen,
				}

				app.baskets.EXPECT().Load(gomock.Any(), gomock.Eq("basket-id")).Times(1).Return(basket, nil)

				app.baskets.EXPECT().Save(gomock.Any(), gomock.AssignableToTypeOf(&domain.Basket{})).Times(1).Return(fmt.Errorf("save failed"))

				app.publisher.EXPECT().Publish(gomock.Any(), gomock.Any()).Times(0)
			},
			wantErr: true,
		},
		{
			name: "PublishFailed",
			args: args{
				ctx:    context.Background(),
				cancel: CancelBasket{ID: "basket-id"},
			},
			buildStubs: func(app MockApplication) {
				basket := &domain.Basket{
					Aggregate:  es.NewAggregate("basket-id", domain.BasketAggregate),
					CustomerID: "customer-id",
					PaymentID:  "payment-id",
					Items:      make(map[string]domain.Item),
					Status:     domain.BasketIsOpen,
				}

				app.baskets.EXPECT().Load(gomock.Any(), gomock.Eq("basket-id")).Times(1).Return(basket, nil)

				app.baskets.EXPECT().Save(gomock.Any(), gomock.AssignableToTypeOf(&domain.Basket{})).Times(1).Return(nil)

				app.publisher.EXPECT().Publish(gomock.Any(), gomock.Any()).Times(1).Return(fmt.Errorf("publish failed"))

			},
			wantErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			m := MockApplication{
				baskets:   domain.NewMockBasketRepository(ctrl),
				stores:    domain.NewMockStoreRepository(ctrl),
				products:  domain.NewMockProductRepository(ctrl),
				publisher: ddd.NewMockEventPublisher[ddd.Event](ctrl),
			}

			app := New(m.baskets, m.stores, m.products, m.publisher)

			tc.buildStubs(m)

			err := app.CancelBasket(tc.args.ctx, tc.args.cancel)
			assert.Equal(t, err != nil, tc.wantErr)
		})
	}
}

func TestApplication_CheckoutBasket(t *testing.T) {
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

	type MockApplication struct {
		baskets   *domain.MockBasketRepository
		stores    *domain.MockStoreRepository
		products  *domain.MockProductRepository
		publisher *ddd.MockEventPublisher[ddd.Event]
	}

	type args struct {
		ctx      context.Context
		checkout CheckoutBasket
	}

	testCases := []struct {
		name       string
		args       args
		buildStubs func(app MockApplication)
		wantErr    bool
	}{
		{
			name: "OK",
			args: args{
				ctx: context.Background(),
				checkout: CheckoutBasket{
					ID:        "basket-id",
					PaymentID: "payment-id",
				},
			},
			buildStubs: func(app MockApplication) {
				basket := &domain.Basket{
					Aggregate:  es.NewAggregate("basket-id", domain.BasketAggregate),
					CustomerID: "customer-id",
					PaymentID:  "payment-id",
					Items:      map[string]domain.Item{product.ID: item},
					Status:     domain.BasketIsOpen,
				}

				app.baskets.EXPECT().Load(gomock.Any(), gomock.Eq("basket-id")).Times(1).Return(basket, nil)

				app.baskets.EXPECT().Save(gomock.Any(), gomock.AssignableToTypeOf(&domain.Basket{})).Times(1).Return(nil)

				app.publisher.EXPECT().Publish(gomock.Any(), gomock.Any()).Times(1).Return(nil)

			},
			wantErr: false,
		},
		{
			name: "NoBasket",
			args: args{
				ctx: context.Background(),
				checkout: CheckoutBasket{
					ID:        "basket-id",
					PaymentID: "payment-id",
				},
			},
			buildStubs: func(app MockApplication) {
				app.baskets.EXPECT().Load(gomock.Any(), gomock.Eq("basket-id")).Times(1).Return(nil, fmt.Errorf("no basket"))

				app.baskets.EXPECT().Save(gomock.Any(), gomock.Any()).Times(0)

				app.publisher.EXPECT().Publish(gomock.Any(), gomock.Any()).Times(0)
			},
			wantErr: true,
		},
		{
			name: "SaveFailed",
			args: args{
				ctx: context.Background(),
				checkout: CheckoutBasket{
					ID:        "basket-id",
					PaymentID: "payment-id",
				},
			},
			buildStubs: func(app MockApplication) {
				basket := &domain.Basket{
					Aggregate:  es.NewAggregate("basket-id", domain.BasketAggregate),
					CustomerID: "customer-id",
					PaymentID:  "payment-id",
					Items:      map[string]domain.Item{product.ID: item},
					Status:     domain.BasketIsOpen,
				}

				app.baskets.EXPECT().Load(gomock.Any(), gomock.Eq("basket-id")).Times(1).Return(basket, nil)

				app.baskets.EXPECT().Save(gomock.Any(), gomock.AssignableToTypeOf(&domain.Basket{})).Times(1).Return(fmt.Errorf("save failed"))

				app.publisher.EXPECT().Publish(gomock.Any(), gomock.Any()).Times(0)
			},
			wantErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			m := MockApplication{
				baskets:   domain.NewMockBasketRepository(ctrl),
				stores:    domain.NewMockStoreRepository(ctrl),
				products:  domain.NewMockProductRepository(ctrl),
				publisher: ddd.NewMockEventPublisher[ddd.Event](ctrl),
			}

			app := New(m.baskets, m.stores, m.products, m.publisher)

			tc.buildStubs(m)

			err := app.CheckoutBasket(tc.args.ctx, tc.args.checkout)
			assert.Equal(t, err != nil, tc.wantErr)
		})
	}
}

func TestApplication_GetBasket(t *testing.T) {
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

	type MockApplication struct {
		baskets   *domain.MockBasketRepository
		stores    *domain.MockStoreRepository
		products  *domain.MockProductRepository
		publisher *ddd.MockEventPublisher[ddd.Event]
	}

	type args struct {
		ctx context.Context
		get GetBasket
	}

	testCases := []struct {
		name       string
		args       args
		buildStubs func(app MockApplication)
		want       *domain.Basket
		wantErr    bool
	}{
		{
			name: "OK",
			args: args{
				ctx: context.Background(),
				get: GetBasket{ID: "basket-id"},
			},
			buildStubs: func(app MockApplication) {
				basket := &domain.Basket{
					Aggregate:  es.NewAggregate("basket-id", domain.BasketAggregate),
					CustomerID: "customer-id",
					PaymentID:  "payment-id",
					Items:      map[string]domain.Item{product.ID: item},
					Status:     domain.BasketIsOpen,
				}

				app.baskets.EXPECT().Load(gomock.Any(), gomock.Eq("basket-id")).Times(1).Return(basket, nil)
			},
			want: &domain.Basket{
				Aggregate:  es.NewAggregate("basket-id", domain.BasketAggregate),
				CustomerID: "customer-id",
				PaymentID:  "payment-id",
				Items:      map[string]domain.Item{product.ID: item},
				Status:     domain.BasketIsOpen,
			},
			wantErr: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			m := MockApplication{
				baskets:   domain.NewMockBasketRepository(ctrl),
				stores:    domain.NewMockStoreRepository(ctrl),
				products:  domain.NewMockProductRepository(ctrl),
				publisher: ddd.NewMockEventPublisher[ddd.Event](ctrl),
			}

			app := New(m.baskets, m.stores, m.products, m.publisher)

			tc.buildStubs(m)

			basket, err := app.GetBasket(tc.args.ctx, tc.args.get)
			assert.Equal(t, err != nil, tc.wantErr)

			assert.Equal(t, tc.want.ID(), basket.ID())
			assert.Equal(t, tc.want.CustomerID, basket.CustomerID)
			assert.Equal(t, tc.want.PaymentID, basket.PaymentID)
			assert.Equal(t, tc.want.Items, basket.Items)
			assert.Equal(t, tc.want.Status, basket.Status)
		})
	}
}

func TestApplication_RemoveItem(t *testing.T) {
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

	type MockApplication struct {
		baskets   *domain.MockBasketRepository
		stores    *domain.MockStoreRepository
		products  *domain.MockProductRepository
		publisher *ddd.MockEventPublisher[ddd.Event]
	}

	type args struct {
		ctx    context.Context
		remove RemoveItem
	}

	testCases := []struct {
		name       string
		args       args
		buildStubs func(app MockApplication)
		wantErr    bool
	}{
		{
			name: "OK",
			args: args{
				ctx: context.Background(),
				remove: RemoveItem{
					ID:        "basket-id",
					ProductID: product.ID,
					Quantity:  1,
				},
			},
			buildStubs: func(app MockApplication) {
				app.products.EXPECT().Find(gomock.Any(), gomock.Eq("product-id")).Times(1).Return(product, nil)

				basket := &domain.Basket{
					Aggregate:  es.NewAggregate("basket-id", domain.BasketAggregate),
					CustomerID: "customer-id",
					PaymentID:  "payment-id",
					Items:      map[string]domain.Item{product.ID: item},
					Status:     domain.BasketIsOpen,
				}

				app.baskets.EXPECT().Load(gomock.Any(), gomock.Eq("basket-id")).Times(1).Return(basket, nil)

				app.baskets.EXPECT().Save(gomock.Any(), gomock.AssignableToTypeOf(&domain.Basket{})).Times(1).Return(nil)
			},
			wantErr: false,
		},
		{
			name: "NoProduct",
			args: args{
				ctx: context.Background(),
				remove: RemoveItem{
					ID:        "basket-id",
					ProductID: product.ID,
					Quantity:  1,
				},
			},
			buildStubs: func(app MockApplication) {
				app.products.EXPECT().Find(gomock.Any(), gomock.Eq("product-id")).Times(1).Return(nil, fmt.Errorf("no product"))

				app.baskets.EXPECT().Load(gomock.Any(), gomock.Any()).Times(0)

				app.baskets.EXPECT().Save(gomock.Any(), gomock.Any()).Times(0)

			},
			wantErr: true,
		},
		{
			name: "NoBasket",
			args: args{
				ctx: context.Background(),
				remove: RemoveItem{
					ID:        "basket-id",
					ProductID: product.ID,
					Quantity:  1,
				},
			},
			buildStubs: func(app MockApplication) {
				app.products.EXPECT().Find(gomock.Any(), gomock.Eq("product-id")).Times(1).Return(product, nil)

				app.baskets.EXPECT().Load(gomock.Any(), gomock.Eq("basket-id")).Times(1).Return(nil, fmt.Errorf("no basket"))

				app.baskets.EXPECT().Save(gomock.Any(), gomock.Any()).Times(0)
			},
			wantErr: true,
		},
		{
			name: "SaveFailed",
			args: args{
				ctx: context.Background(),
				remove: RemoveItem{
					ID:        "basket-id",
					ProductID: product.ID,
					Quantity:  1,
				},
			},
			buildStubs: func(app MockApplication) {
				app.products.EXPECT().Find(gomock.Any(), gomock.Eq("product-id")).Times(1).Return(product, nil)

				basket := &domain.Basket{
					Aggregate:  es.NewAggregate("basket-id", domain.BasketAggregate),
					CustomerID: "customer-id",
					PaymentID:  "payment-id",
					Items:      map[string]domain.Item{product.ID: item},
					Status:     domain.BasketIsOpen,
				}

				app.baskets.EXPECT().Load(gomock.Any(), gomock.Eq("basket-id")).Times(1).Return(basket, nil)

				app.baskets.EXPECT().Save(gomock.Any(), gomock.AssignableToTypeOf(&domain.Basket{})).Times(1).Return(fmt.Errorf("save failed"))
			},
			wantErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			m := MockApplication{
				baskets:   domain.NewMockBasketRepository(ctrl),
				stores:    domain.NewMockStoreRepository(ctrl),
				products:  domain.NewMockProductRepository(ctrl),
				publisher: ddd.NewMockEventPublisher[ddd.Event](ctrl),
			}

			app := New(m.baskets, m.stores, m.products, m.publisher)

			tc.buildStubs(m)

			err := app.RemoveItem(tc.args.ctx, tc.args.remove)
			assert.Equal(t, err != nil, tc.wantErr)
		})
	}
}

func TestApplication_StartBasket(t *testing.T) {
	type MockApplication struct {
		baskets   *domain.MockBasketRepository
		stores    *domain.MockStoreRepository
		products  *domain.MockProductRepository
		publisher *ddd.MockEventPublisher[ddd.Event]
	}

	type args struct {
		ctx   context.Context
		start StartBasket
	}

	testCases := []struct {
		name       string
		args       args
		buildStubs func(app MockApplication)
		wantErr    bool
	}{
		{
			name: "OK",
			args: args{
				ctx: context.Background(),
				start: StartBasket{
					ID:         "basket-id",
					CustomerID: "customer-id",
				},
			},
			buildStubs: func(app MockApplication) {
				basket := &domain.Basket{
					Aggregate:  es.NewAggregate("basket-id", domain.BasketAggregate),
					CustomerID: "customer-id",
					PaymentID:  "payment-id",
					Items:      make(map[string]domain.Item),
				}

				app.baskets.EXPECT().Load(gomock.Any(), gomock.Eq("basket-id")).Times(1).Return(basket, nil)

				app.baskets.EXPECT().Save(gomock.Any(), gomock.AssignableToTypeOf(&domain.Basket{})).Times(1).Return(nil)

				app.publisher.EXPECT().Publish(gomock.Any(), gomock.Any()).Times(1).Return(nil)
			},
			wantErr: false,
		},
		{
			name: "NoBasket",
			args: args{
				ctx: context.Background(),
				start: StartBasket{
					ID:         "basket-id",
					CustomerID: "customer-id",
				},
			},
			buildStubs: func(app MockApplication) {
				app.baskets.EXPECT().Load(gomock.Any(), gomock.Eq("basket-id")).Times(1).Return(nil, fmt.Errorf("no basket"))

				app.baskets.EXPECT().Save(gomock.Any(), gomock.Any()).Times(0)

				app.publisher.EXPECT().Publish(gomock.Any(), gomock.Any()).Times(0)
			},
			wantErr: true,
		},
		{
			name: "SaveFailed",
			args: args{
				ctx: context.Background(),
				start: StartBasket{
					ID:         "basket-id",
					CustomerID: "customer-id",
				},
			},
			buildStubs: func(app MockApplication) {
				basket := &domain.Basket{
					Aggregate:  es.NewAggregate("basket-id", domain.BasketAggregate),
					CustomerID: "customer-id",
					PaymentID:  "payment-id",
					Items:      make(map[string]domain.Item),
				}

				app.baskets.EXPECT().Load(gomock.Any(), gomock.Eq("basket-id")).Times(1).Return(basket, nil)

				app.baskets.EXPECT().Save(gomock.Any(), gomock.AssignableToTypeOf(&domain.Basket{})).Times(1).Return(fmt.Errorf("save failed"))

				app.publisher.EXPECT().Publish(gomock.Any(), gomock.Any()).Times(0)
			},
			wantErr: true,
		},
		{
			name: "PublishFailed",
			args: args{
				ctx: context.Background(),
				start: StartBasket{
					ID:         "basket-id",
					CustomerID: "customer-id",
				},
			},
			buildStubs: func(app MockApplication) {
				basket := &domain.Basket{
					Aggregate:  es.NewAggregate("basket-id", domain.BasketAggregate),
					CustomerID: "customer-id",
					PaymentID:  "payment-id",
					Items:      make(map[string]domain.Item),
				}

				app.baskets.EXPECT().Load(gomock.Any(), gomock.Eq("basket-id")).Times(1).Return(basket, nil)

				app.baskets.EXPECT().Save(gomock.Any(), gomock.AssignableToTypeOf(&domain.Basket{})).Times(1).Return(nil)

				app.publisher.EXPECT().Publish(gomock.Any(), gomock.Any()).Times(1).Return(fmt.Errorf("publish failed"))

			},
			wantErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			m := MockApplication{
				baskets:   domain.NewMockBasketRepository(ctrl),
				stores:    domain.NewMockStoreRepository(ctrl),
				products:  domain.NewMockProductRepository(ctrl),
				publisher: ddd.NewMockEventPublisher[ddd.Event](ctrl),
			}

			app := New(m.baskets, m.stores, m.products, m.publisher)

			tc.buildStubs(m)

			err := app.StartBasket(tc.args.ctx, tc.args.start)
			assert.Equal(t, err != nil, tc.wantErr)
		})
	}
}
