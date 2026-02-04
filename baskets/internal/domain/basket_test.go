package domain

import (
	"mall/internal/ddd"
	"mall/internal/es"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestBasket_AddItem(t *testing.T) {
	store := &Store{
		ID:   "store-id",
		Name: "store-name",
	}

	product := &Product{
		ID:      "product-id",
		StoreID: "store-id",
		Name:    "product-name",
		Price:   10.00,
	}

	type basketState struct {
		CustomerID string
		PaymentID  string
		Items      map[string]Item
		Status     BasketStatus
	}

	type args struct {
		store    *Store
		product  *Product
		quantity int
	}

	testCases := []struct {
		name        string
		basketState basketState
		args        args
		buildStubs  func(a *es.MockAggregate)
		wantErr     bool
	}{
		{
			name: "OK",
			basketState: basketState{
				Items:  make(map[string]Item),
				Status: BasketIsOpen,
			},
			args: args{
				store:    store,
				product:  product,
				quantity: 1,
			},
			buildStubs: func(a *es.MockAggregate) {
				basketAddedEvent := &BasketItemAdded{
					Item: Item{
						StoreID:      store.ID,
						ProductID:    product.ID,
						StoreName:    store.Name,
						ProductName:  product.Name,
						ProductPrice: product.Price,
						Quantity:     1,
					},
				}

				a.EXPECT().AddEvent(gomock.Eq(BasketItemAddedEvent), basketAddedEvent).Times(1)
			},
			wantErr: false,
		},
		{
			name: "CanceledOutBasket",
			basketState: basketState{
				Items:  make(map[string]Item),
				Status: BasketIsCanceled,
			},
			args: args{
				store:    store,
				product:  product,
				quantity: 1,
			},
			buildStubs: func(a *es.MockAggregate) {
				a.EXPECT().AddEvent(gomock.Any(), gomock.Any()).Times(0)
			},
			wantErr: true,
		},
		{
			name: "CheckedOutBasket",
			basketState: basketState{
				Items:  make(map[string]Item),
				Status: BasketIsCheckedOut,
			},
			args: args{
				store:    store,
				product:  product,
				quantity: 1,
			},
			buildStubs: func(a *es.MockAggregate) {
				a.EXPECT().AddEvent(gomock.Any(), gomock.Any()).Times(0)
			},
			wantErr: true,
		},
		{
			name: "ZeroQuantity",
			basketState: basketState{
				Items:  make(map[string]Item),
				Status: BasketIsOpen,
			},
			args: args{
				store:    store,
				product:  product,
				quantity: 0,
			},
			buildStubs: func(a *es.MockAggregate) {
				a.EXPECT().AddEvent(gomock.Any(), gomock.Any()).Times(0)
			},
			wantErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			aggregate := es.NewMockAggregate(ctrl)

			b := &Basket{
				Aggregate:  aggregate,
				CustomerID: tc.basketState.CustomerID,
				PaymentID:  tc.basketState.PaymentID,
				Items:      tc.basketState.Items,
				Status:     tc.basketState.Status,
			}

			tc.buildStubs(aggregate)

			err := b.AddItem(tc.args.store, tc.args.product, tc.args.quantity)
			assert.Equal(t, err != nil, tc.wantErr)
		})
	}
}

func TestBasket_ApplyEvent(t *testing.T) {
	store := &Store{
		ID:   "store-id",
		Name: "store-name",
	}

	product := &Product{
		ID:      "product-id",
		StoreID: "store-id",
		Name:    "product-name",
		Price:   10.00,
	}

	product1 := &Product{
		ID:      "product-id1",
		StoreID: "store-id",
		Name:    "product_name1",
		Price:   100.00,
	}

	type basketState struct {
		CustomerID string
		PaymentID  string
		Items      map[string]Item
		Status     BasketStatus
	}

	type args struct {
		event ddd.Event
	}

	testCases := []struct {
		name        string
		basketState basketState
		args        args
		wantState   basketState
		wantErr     bool
	}{
		{
			name: "BasketItemAddedEvent",
			basketState: basketState{
				CustomerID: "customer-id",
				PaymentID:  "payment-id",
				Items:      make(map[string]Item),
				Status:     BasketIsOpen,
			},
			args: args{
				event: ddd.NewEvent(BasketItemAddedEvent, &BasketItemAdded{
					Item: Item{
						StoreID:      store.ID,
						ProductID:    product.ID,
						StoreName:    store.Name,
						ProductName:  product.Name,
						ProductPrice: product.Price,
						Quantity:     1,
					},
				}),
			},
			wantState: basketState{
				CustomerID: "customer-id",
				PaymentID:  "payment-id",
				Items: map[string]Item{
					product.ID: {
						StoreID:      store.ID,
						ProductID:    product.ID,
						StoreName:    store.Name,
						ProductName:  product.Name,
						ProductPrice: product.Price,
						Quantity:     1,
					},
				},
				Status: BasketIsOpen,
			},
			wantErr: false,
		},
		{
			name: "BasketItemAddedEvent.Quantity",
			basketState: basketState{
				CustomerID: "customer-id",
				PaymentID:  "payment-id",
				Items: map[string]Item{
					product.ID: {
						StoreID:      store.ID,
						ProductID:    product.ID,
						StoreName:    store.Name,
						ProductName:  product.Name,
						ProductPrice: product.Price,
						Quantity:     1,
					},
				},
				Status: BasketIsOpen,
			},
			args: args{
				event: ddd.NewEvent(BasketItemAddedEvent, &BasketItemAdded{
					Item: Item{
						StoreID:      store.ID,
						ProductID:    product.ID,
						StoreName:    store.Name,
						ProductName:  product.Name,
						ProductPrice: product.Price,
						Quantity:     1,
					},
				}),
			},
			wantState: basketState{
				CustomerID: "customer-id",
				PaymentID:  "payment-id",
				Items: map[string]Item{
					product.ID: {
						StoreID:      store.ID,
						ProductID:    product.ID,
						StoreName:    store.Name,
						ProductName:  product.Name,
						ProductPrice: product.Price,
						Quantity:     2,
					},
				},
				Status: BasketIsOpen,
			},
			wantErr: false,
		},
		{
			name: "BasketItemAddedEvent.Second",
			basketState: basketState{
				CustomerID: "customer-id",
				PaymentID:  "payment-id",
				Items: map[string]Item{
					product.ID: {
						StoreID:      store.ID,
						ProductID:    product.ID,
						StoreName:    store.Name,
						ProductName:  product.Name,
						ProductPrice: product.Price,
						Quantity:     1,
					},
				},
				Status: BasketIsOpen,
			},
			args: args{
				event: ddd.NewEvent(BasketItemAddedEvent, &BasketItemAdded{
					Item: Item{
						StoreID:      store.ID,
						ProductID:    product1.ID,
						StoreName:    store.Name,
						ProductName:  product1.Name,
						ProductPrice: product1.Price,
						Quantity:     1,
					},
				}),
			},
			wantState: basketState{
				CustomerID: "customer-id",
				PaymentID:  "payment-id",
				Items: map[string]Item{
					product.ID: {
						StoreID:      store.ID,
						ProductID:    product.ID,
						StoreName:    store.Name,
						ProductName:  product.Name,
						ProductPrice: product.Price,
						Quantity:     1,
					},
					product1.ID: {
						StoreID:      store.ID,
						ProductID:    product1.ID,
						StoreName:    store.Name,
						ProductName:  product1.Name,
						ProductPrice: product1.Price,
						Quantity:     1,
					},
				},
				Status: BasketIsOpen,
			},
			wantErr: false,
		},
		{
			name: "BasketCanceledEvent",
			basketState: basketState{
				CustomerID: "customer-id",
				PaymentID:  "payment-id",
				Items:      make(map[string]Item),
				Status:     BasketIsOpen,
			},
			args: args{
				event: ddd.NewEvent(BasketCanceledEvent, &BasketCanceled{}),
			},
			wantState: basketState{
				CustomerID: "customer-id",
				PaymentID:  "payment-id",
				Items:      map[string]Item{},
				Status:     BasketIsCanceled,
			},
			wantErr: false,
		},
		{
			name: "BasketCanceledEvent.Cleared",
			basketState: basketState{
				CustomerID: "customer-id",
				PaymentID:  "payment-id",
				Items: map[string]Item{
					product.ID: {
						StoreID:      store.ID,
						ProductID:    product.ID,
						StoreName:    store.Name,
						ProductName:  product.Name,
						ProductPrice: product.Price,
						Quantity:     1,
					},
				},
				Status: BasketIsOpen,
			},
			args: args{
				event: ddd.NewEvent(BasketCanceledEvent, &BasketCanceled{}),
			},
			wantState: basketState{
				CustomerID: "customer-id",
				PaymentID:  "payment-id",
				Items:      map[string]Item{},
				Status:     BasketIsCanceled,
			},
			wantErr: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			aggregate := es.NewMockAggregate(ctrl)

			b := &Basket{
				Aggregate:  aggregate,
				CustomerID: tc.basketState.CustomerID,
				PaymentID:  tc.basketState.PaymentID,
				Items:      tc.basketState.Items,
				Status:     tc.basketState.Status,
			}

			err := b.ApplyEvent(tc.args.event)
			assert.Equal(t, err != nil, tc.wantErr)

			assert.Equal(t, b.CustomerID, tc.wantState.CustomerID)
			assert.Equal(t, b.PaymentID, tc.wantState.PaymentID)
			assert.Equal(t, b.Items, tc.wantState.Items)
			assert.Equal(t, b.Status, tc.wantState.Status)
		})
	}
}

func TestBasket_ApplySnapshot(t *testing.T) {
	store := &Store{
		ID:   "store-id",
		Name: "store-name",
	}

	product := &Product{
		ID:      "product-id",
		StoreID: "store-id",
		Name:    "product-name",
		Price:   10.00,
	}

	type basketState struct {
		CustomerID string
		PaymentID  string
		Items      map[string]Item
		Status     BasketStatus
	}

	type args struct {
		snapshot es.Snapshot
	}

	item := Item{
		StoreID:      store.ID,
		ProductID:    product.ID,
		StoreName:    store.Name,
		ProductName:  product.Name,
		ProductPrice: product.Price,
		Quantity:     1,
	}

	testCases := []struct {
		name        string
		basketState basketState
		args        args
		wantState   basketState
		wantErr     bool
	}{
		{
			name:        "V1",
			basketState: basketState{},
			args: args{
				snapshot: &BasketV1{
					CustomerID: "customer-id",
					PaymentID:  "payment-id",
					Items: map[string]Item{
						product.ID: item,
					},
					Status: BasketIsOpen,
				},
			},
			wantState: basketState{
				CustomerID: "customer-id",
				PaymentID:  "payment-id",
				Items: map[string]Item{
					product.ID: item,
				},
				Status: BasketIsOpen,
			},
			wantErr: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			aggregate := es.NewMockAggregate(ctrl)

			b := &Basket{
				Aggregate:  aggregate,
				CustomerID: tc.basketState.CustomerID,
				PaymentID:  tc.basketState.PaymentID,
				Items:      tc.basketState.Items,
				Status:     tc.basketState.Status,
			}

			err := b.ApplySnapshot(tc.args.snapshot)
			assert.Equal(t, err != nil, tc.wantErr)

			assert.Equal(t, b.CustomerID, tc.wantState.CustomerID)
			assert.Equal(t, b.PaymentID, tc.wantState.PaymentID)
			assert.Equal(t, b.Items, tc.wantState.Items)
			assert.Equal(t, b.Status, tc.wantState.Status)
		})
	}
}

func TestBasket_Cancel(t *testing.T) {
	type basketState struct {
		CustomerID string
		PaymentID  string
		Items      map[string]Item
		Status     BasketStatus
	}

	testCases := []struct {
		name        string
		basketState basketState
		buildStubs  func(a *es.MockAggregate)
		wantState   ddd.Event
		wantErr     bool
	}{
		{
			name: "OpenBasket",
			basketState: basketState{
				CustomerID: "customer-id",
				PaymentID:  "payment-id",
				Items:      make(map[string]Item),
				Status:     BasketIsOpen,
			},
			buildStubs: func(a *es.MockAggregate) {
				a.EXPECT().AddEvent(gomock.Eq(BasketCanceledEvent), gomock.AssignableToTypeOf(&BasketCanceled{})).Times(1)
			},
			wantState: ddd.NewEvent(BasketCanceledEvent, &Basket{
				CustomerID: "customer-id",
				PaymentID:  "payment-id",
				Items:      make(map[string]Item),
				Status:     BasketIsCanceled,
			}),
			wantErr: false,
		},
		{
			name: "CancelBasket",
			basketState: basketState{
				CustomerID: "customer-id",
				PaymentID:  "payment-id",
				Items:      make(map[string]Item),
				Status:     BasketIsCanceled,
			},
			buildStubs: func(a *es.MockAggregate) {
				a.EXPECT().AddEvent(gomock.Any(), gomock.Any()).Times(0)
			},
			wantState: nil,
			wantErr:   true,
		},
		{
			name: "CheckedOutBasket",
			basketState: basketState{
				CustomerID: "customer-id",
				PaymentID:  "payment-id",
				Items:      make(map[string]Item),
				Status:     BasketIsCheckedOut,
			},
			buildStubs: func(a *es.MockAggregate) {
				a.EXPECT().AddEvent(gomock.Any(), gomock.Any()).Times(0)
			},
			wantState: nil,
			wantErr:   true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			aggregate := es.NewMockAggregate(ctrl)

			b := &Basket{
				Aggregate:  aggregate,
				CustomerID: tc.basketState.CustomerID,
				PaymentID:  tc.basketState.PaymentID,
				Items:      tc.basketState.Items,
				Status:     tc.basketState.Status,
			}

			tc.buildStubs(aggregate)

			event, err := b.Cancel()
			assert.Equal(t, err != nil, tc.wantErr)

			if tc.wantState != nil {
				assert.Equal(t, event.EventName(), tc.wantState.EventName())
				assert.IsType(t, event.Payload(), tc.wantState.Payload())
				assert.Equal(t, event.Metadata(), tc.wantState.Metadata())
			} else {
				assert.Nil(t, event)
			}
		})
	}
}

func TestBasket_Checkout(t *testing.T) {
	store := &Store{
		ID:   "store-id",
		Name: "store-name",
	}

	product := &Product{
		ID:      "product-id",
		StoreID: "store-id",
		Name:    "product-name",
		Price:   10.00,
	}

	item := Item{
		StoreID:      store.ID,
		ProductID:    product.ID,
		StoreName:    store.Name,
		ProductName:  product.Name,
		ProductPrice: product.Price,
		Quantity:     1,
	}

	type basketState struct {
		CustomerID string
		PaymentID  string
		Items      map[string]Item
		Status     BasketStatus
	}

	type args struct {
		paymentID string
	}

	testCases := []struct {
		name        string
		basketState basketState
		args        args
		buildStubs  func(a *es.MockAggregate)
		wantState   ddd.Event
		wantErr     bool
	}{
		{
			name: "OpenBasket",
			basketState: basketState{
				CustomerID: "customer-id",
				PaymentID:  "payment-id",
				Items: map[string]Item{
					product.ID: item,
				},
				Status: BasketIsOpen,
			},
			args: args{paymentID: "payment-id"},
			buildStubs: func(a *es.MockAggregate) {
				a.EXPECT().AddEvent(gomock.Eq(BasketCheckedOutEvent), &BasketCheckedOut{PaymentID: "payment-id"}).Times(1)
			},
			wantState: ddd.NewEvent(BasketCheckedOutEvent, &Basket{
				CustomerID: "customer-id",
				PaymentID:  "payment-id",
				Items:      make(map[string]Item),
				Status:     BasketIsCheckedOut,
			}),
			wantErr: false,
		},
		{
			name: "OpenBasket.NoItems",
			basketState: basketState{
				CustomerID: "customer-id",
				PaymentID:  "payment-id",
				Items:      map[string]Item{},
				Status:     BasketIsOpen,
			},
			args: args{paymentID: "payment-id"},
			buildStubs: func(a *es.MockAggregate) {
				a.EXPECT().AddEvent(gomock.Any(), gomock.Any()).Times(0)
			},
			wantState: nil,
			wantErr:   true,
		},
		{
			name: "CanceledBasket",
			basketState: basketState{
				CustomerID: "customer-id",
				PaymentID:  "payment-id",
				Items:      map[string]Item{},
				Status:     BasketIsCanceled,
			},
			args: args{paymentID: "payment-id"},
			buildStubs: func(a *es.MockAggregate) {
				a.EXPECT().AddEvent(gomock.Any(), gomock.Any()).Times(0)
			},
			wantState: nil,
			wantErr:   true,
		},
		{
			name: "CheckedOutBasket",
			basketState: basketState{
				CustomerID: "customer-id",
				PaymentID:  "payment-id",
				Items:      map[string]Item{},
				Status:     BasketIsCheckedOut,
			},
			args: args{paymentID: "payment-id"},
			buildStubs: func(a *es.MockAggregate) {
				a.EXPECT().AddEvent(gomock.Any(), gomock.Any()).Times(0)
			},
			wantState: nil,
			wantErr:   true,
		},
		{
			name: "NoPaymentId",
			basketState: basketState{
				CustomerID: "customer-id",
				PaymentID:  "",
				Items:      map[string]Item{},
				Status:     BasketIsCheckedOut,
			},
			args: args{paymentID: ""},
			buildStubs: func(a *es.MockAggregate) {
				a.EXPECT().AddEvent(gomock.Any(), gomock.Any()).Times(0)
			},
			wantState: nil,
			wantErr:   true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			aggregate := es.NewMockAggregate(ctrl)

			b := &Basket{
				Aggregate:  aggregate,
				CustomerID: tc.basketState.CustomerID,
				PaymentID:  tc.basketState.PaymentID,
				Items:      tc.basketState.Items,
				Status:     tc.basketState.Status,
			}

			tc.buildStubs(aggregate)

			event, err := b.Checkout(tc.args.paymentID)
			assert.Equal(t, err != nil, tc.wantErr)

			if tc.wantState != nil {
				assert.Equal(t, event.EventName(), tc.wantState.EventName())
				assert.IsType(t, event.Payload(), tc.wantState.Payload())
				assert.Equal(t, event.Metadata(), tc.wantState.Metadata())
			} else {
				assert.Nil(t, event)
			}
		})
	}
}

func TestBasket_RemoveItem(t *testing.T) {
	store := &Store{
		ID:   "store-id",
		Name: "store-name",
	}

	product := &Product{
		ID:      "product-id",
		StoreID: "store-id",
		Name:    "product-name",
		Price:   10.00,
	}

	item := Item{
		StoreID:      store.ID,
		ProductID:    product.ID,
		StoreName:    store.Name,
		ProductName:  product.Name,
		ProductPrice: product.Price,
		Quantity:     10,
	}

	type basketState struct {
		CustomerID string
		PaymentID  string
		Items      map[string]Item
		Status     BasketStatus
	}

	type args struct {
		product  *Product
		quantity int
	}

	testCases := []struct {
		name        string
		basketState basketState
		args        args
		buildStubs  func(a *es.MockAggregate)
		wantErr     bool
	}{
		{
			name: "OpenBasket",
			basketState: basketState{
				CustomerID: "customer-id",
				PaymentID:  "payment-id",
				Items: map[string]Item{
					product.ID: item,
				},
				Status: BasketIsOpen,
			},
			args: args{
				product:  product,
				quantity: 1,
			},
			buildStubs: func(a *es.MockAggregate) {
				a.EXPECT().AddEvent(gomock.Eq(BasketItemRemovedEvent), &BasketItemRemoved{
					ProductID: product.ID,
					Quantity:  1,
				}).Times(1)
			},
			wantErr: false,
		},
		{
			name: "OpenBasket.NoItems",
			basketState: basketState{
				CustomerID: "customer-id",
				PaymentID:  "payment-id",
				Items:      map[string]Item{},
				Status:     BasketIsOpen,
			},
			args: args{
				product:  product,
				quantity: 1,
			},
			buildStubs: func(a *es.MockAggregate) {
				a.EXPECT().AddEvent(gomock.Any(), gomock.Any()).Times(0)
			},
			wantErr: true,
		},
		{
			name: "CanceledBasket",
			basketState: basketState{
				CustomerID: "customer-id",
				PaymentID:  "payment-id",
				Items:      map[string]Item{},
				Status:     BasketIsCanceled,
			},
			args: args{
				product:  product,
				quantity: 1,
			},
			buildStubs: func(a *es.MockAggregate) {
				a.EXPECT().AddEvent(gomock.Any(), gomock.Any()).Times(0)
			},
			wantErr: true,
		},
		{
			name: "CheckedOutBasket",
			basketState: basketState{
				CustomerID: "customer-id",
				PaymentID:  "payment-id",
				Items:      map[string]Item{},
				Status:     BasketIsCheckedOut,
			},
			args: args{
				product:  product,
				quantity: 1,
			},
			buildStubs: func(a *es.MockAggregate) {
				a.EXPECT().AddEvent(gomock.Any(), gomock.Any()).Times(0)
			},
			wantErr: true,
		},
		{
			name: "ZeroQuantity",
			basketState: basketState{
				CustomerID: "customer-id",
				PaymentID:  "payment-id",
				Items: map[string]Item{
					product.ID: item,
				},
				Status: BasketIsOpen,
			},
			args: args{
				product:  product,
				quantity: 0,
			},
			buildStubs: func(a *es.MockAggregate) {
				a.EXPECT().AddEvent(gomock.Any(), gomock.Any()).Times(0)
			},
			wantErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			aggregate := es.NewMockAggregate(ctrl)

			b := &Basket{
				Aggregate:  aggregate,
				CustomerID: tc.basketState.CustomerID,
				PaymentID:  tc.basketState.PaymentID,
				Items:      tc.basketState.Items,
				Status:     tc.basketState.Status,
			}

			tc.buildStubs(aggregate)

			err := b.RemoveItem(tc.args.product, tc.args.quantity)
			assert.Equal(t, err != nil, tc.wantErr)
		})
	}
}

func TestBasket_Start(t *testing.T) {
	type basketState struct {
		CustomerID string
		PaymentID  string
		Items      map[string]Item
		Status     BasketStatus
	}

	type args struct {
		customerID string
	}

	testCases := []struct {
		name        string
		basketState basketState
		args        args
		buildStubs  func(a *es.MockAggregate)
		wantState   ddd.Event
		wantErr     bool
	}{
		{
			name:        "OK",
			basketState: basketState{},
			args:        args{customerID: "customer-id"},
			buildStubs: func(a *es.MockAggregate) {
				a.EXPECT().AddEvent(gomock.Eq(BasketStartedEvent), &BasketStarted{CustomerID: "customer-id"}).Times(1)
			},
			wantState: ddd.NewEvent(BasketStartedEvent, &Basket{
				CustomerID: "customer-id",
				PaymentID:  "",
				Items:      make(map[string]Item),
				Status:     BasketIsOpen,
			}),
			wantErr: false,
		},
		{
			name:        "OpenBasket",
			basketState: basketState{Status: BasketIsOpen},
			args:        args{customerID: "customer-id"},
			buildStubs: func(a *es.MockAggregate) {
				a.EXPECT().AddEvent(gomock.Any(), gomock.Any()).Times(0)
			},
			wantState: nil,
			wantErr:   true,
		},
		{
			name:        "EmptyCustomerID",
			basketState: basketState{},
			args:        args{customerID: ""},
			buildStubs: func(a *es.MockAggregate) {
				a.EXPECT().AddEvent(gomock.Any(), gomock.Any()).Times(0)
			},
			wantState: nil,
			wantErr:   true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			aggregate := es.NewMockAggregate(ctrl)

			b := &Basket{
				Aggregate:  aggregate,
				CustomerID: tc.basketState.CustomerID,
				PaymentID:  tc.basketState.PaymentID,
				Items:      tc.basketState.Items,
				Status:     tc.basketState.Status,
			}

			tc.buildStubs(aggregate)

			event, err := b.Start(tc.args.customerID)
			assert.Equal(t, err != nil, tc.wantErr)

			if tc.wantState != nil {
				assert.Equal(t, event.EventName(), tc.wantState.EventName())
				assert.IsType(t, event.Payload(), tc.wantState.Payload())
				assert.Equal(t, event.Metadata(), tc.wantState.Metadata())
			} else {
				assert.Nil(t, event)
			}
		})
	}
}

func TestBasket_ToSnapshot(t *testing.T) {
	store := &Store{
		ID:   "store-id",
		Name: "store-name",
	}

	product := &Product{
		ID:      "product-id",
		StoreID: "store-id",
		Name:    "product-name",
		Price:   10.00,
	}

	item := Item{
		StoreID:      store.ID,
		ProductID:    product.ID,
		StoreName:    store.Name,
		ProductName:  product.Name,
		ProductPrice: product.Price,
		Quantity:     1,
	}

	type basketState struct {
		CustomerID string
		PaymentID  string
		Items      map[string]Item
		Status     BasketStatus
	}

	testCases := []struct {
		name        string
		basketState basketState
		wantState   es.Snapshot
	}{
		{
			name: "V1",
			basketState: basketState{
				CustomerID: "customer-id",
				PaymentID:  "payment-id",
				Items: map[string]Item{
					product.ID: item,
				},
				Status: BasketIsOpen,
			},
			wantState: &BasketV1{
				CustomerID: "customer-id",
				PaymentID:  "payment-id",
				Items: map[string]Item{
					product.ID: item,
				},
				Status: BasketIsOpen,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			aggregate := es.NewMockAggregate(ctrl)

			b := &Basket{
				Aggregate:  aggregate,
				CustomerID: tc.basketState.CustomerID,
				PaymentID:  tc.basketState.PaymentID,
				Items:      tc.basketState.Items,
				Status:     tc.basketState.Status,
			}

			snapshot := b.ToSnapshot()
			assert.True(t, reflect.DeepEqual(snapshot, tc.wantState))
		})
	}
}

func TestNewBasket(t *testing.T) {
	type args struct {
		id string
	}

	testCases := []struct {
		name string
		args args
		want *Basket
	}{
		{
			name: "Basket",
			args: args{id: "basket-id"},
			want: &Basket{
				Aggregate: es.NewAggregate("basket-id", BasketAggregate),
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			basket := NewBasket(tc.args.id)

			assert.Equal(t, basket.ID(), tc.want.ID())
			assert.Equal(t, basket.AggregateName(), tc.want.AggregateName())
		})
	}
}
