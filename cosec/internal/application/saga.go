package application

import (
	"context"
	"mall/cosec/internal/domain"
	"mall/customers/customerspb"
	"mall/depot/depotpb"
	"mall/internal/ddd"
	"mall/internal/sec"
	"mall/ordering/orderingpb"
	"mall/payments/paymentspb"
)

type createOrderSaga struct {
	sec.Saga[*domain.CreateOrderData]
}

func NewCreateOrderSaga() sec.Saga[*domain.CreateOrderData] {
	saga := createOrderSaga{
		Saga: sec.NewSaga[*domain.CreateOrderData](
			domain.CreateOrderSagaName,
			domain.CreateOrderReplyChannel,
		),
	}

	// 0. +RejectOrder
	saga.AddStep().
		Compensation(saga.rejectOrder)

	// 1. +AuthorizeCustomer
	saga.AddStep().
		Action(saga.authorizeCustomer)

	// 2. +CreateShoppingList, -CancelShoppingList
	saga.AddStep().
		Action(saga.createShoppingList).
		OnActionReply(depotpb.CreatedShoppingListReply, saga.onCreatedShoppingListReply).
		Compensation(saga.cancelShoppingList)

	// 3. +ConfirmPayment
	saga.AddStep().
		Action(saga.confirmPayment)

	// 4. +InitiateShopping
	saga.AddStep().
		Action(saga.initiateShopping)

	// 5. +ApproveOrder
	saga.AddStep().
		Action(saga.approveOrder)

	return saga
}

func (s createOrderSaga) rejectOrder(ctx context.Context, data *domain.CreateOrderData) (string, ddd.Command, error) {
	command := ddd.NewCommand(orderingpb.RejectOrderCommand, &orderingpb.RejectOrder{
		Id: data.OrderID,
	})

	return orderingpb.CommandChannel, command, nil
}

func (s createOrderSaga) authorizeCustomer(ctx context.Context, data *domain.CreateOrderData) (string, ddd.Command, error) {
	command := ddd.NewCommand(customerspb.AuthorizeCustomerCommand, &customerspb.AuthorizeCustomer{
		Id: data.CustomerID,
	})

	return customerspb.CommandChannel, command, nil
}

func (s createOrderSaga) createShoppingList(ctx context.Context, data *domain.CreateOrderData) (string, ddd.Command, error) {
	items := make([]*depotpb.CreateShoppingList_Item, len(data.Items))

	for i, item := range data.Items {
		items[i] = &depotpb.CreateShoppingList_Item{
			ProductId: item.ProductID,
			StoreId:   item.StoreID,
			Quantity:  int32(item.Quantity),
		}
	}

	command := ddd.NewCommand(depotpb.CreateShoppingListCommand, &depotpb.CreateShoppingList{
		OrderId: data.OrderID,
		Items:   items,
	})

	return depotpb.CommandChannel, command, nil
}

func (s createOrderSaga) onCreatedShoppingListReply(ctx context.Context, data *domain.CreateOrderData, reply ddd.Reply) error {
	payload := reply.Payload().(*depotpb.CreatedShoppingList)

	data.ShoppingID = payload.GetId()

	return nil
}

func (s createOrderSaga) cancelShoppingList(ctx context.Context, data *domain.CreateOrderData) (string, ddd.Command, error) {
	command := ddd.NewCommand(depotpb.CancelShoppingListCommand, &depotpb.CancelShoppingList{
		Id: data.ShoppingID,
	})

	return depotpb.CommandChannel, command, nil
}

func (s createOrderSaga) confirmPayment(ctx context.Context, data *domain.CreateOrderData) (string, ddd.Command, error) {
	command := ddd.NewCommand(paymentspb.ConfirmPaymentCommand, &paymentspb.ConfirmPayment{
		Id:     data.PaymentID,
		Amount: data.Total,
	})

	return paymentspb.CommandChannel, command, nil
}

func (s createOrderSaga) initiateShopping(ctx context.Context, data *domain.CreateOrderData) (string, ddd.Command, error) {
	command := ddd.NewCommand(depotpb.InitiateShoppingCommand, &depotpb.InitiateShopping{
		Id: data.ShoppingID,
	})

	return depotpb.CommandChannel, command, nil
}

func (s createOrderSaga) approveOrder(ctx context.Context, data *domain.CreateOrderData) (string, ddd.Command, error) {
	command := ddd.NewCommand(orderingpb.ApproveOrderCommand, &orderingpb.ApproveOrder{
		Id:         data.OrderID,
		ShoppingId: data.ShoppingID,
	})

	return orderingpb.CommandChannel, command, nil
}
