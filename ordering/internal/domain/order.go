package domain

import (
	"errors"
	"mall/internal/ddd"
)

var (
	ErrOrderHasNoItems         = errors.New("the order has no items")
	ErrOrderCannotBeCancelled  = errors.New("the order cannot be cancelled")
	ErrOrderCannotBeReady      = errors.New("the order cannot be readied")
	ErrOrderCannotBeCompleted  = errors.New("the order cannot be completed")
	ErrCustomerIDCannotBeBlank = errors.New("the customer id cannot be blank")
	ErrPaymentIDCannotBeBlank  = errors.New("the payment id cannot be blank")
	ErrInvoiceIDCannotBeBlank  = errors.New("the invoice id cannot be blank")
)

type Order struct {
	ddd.AggregateBase
	CustomerID string
	PaymentID  string
	InvoiceID  string
	ShoppingID string
	Items      []*Item
	Status     OrderStatus
}

func CreateOrder(id, customerID, paymentID string, items []*Item) (*Order, error) {
	if len(items) == 0 {
		return nil, ErrOrderHasNoItems
	}

	if customerID == "" {
		return nil, ErrCustomerIDCannotBeBlank
	}

	if paymentID == "" {
		return nil, ErrPaymentIDCannotBeBlank
	}

	order := &Order{
		AggregateBase: ddd.AggregateBase{ID: id},
		CustomerID:    customerID,
		PaymentID:     paymentID,
		Items:         items,
		Status:        OrderIsPending,
	}

	order.AddEvent(&OrderCreated{
		Order: order,
	})

	return order, nil
}

func (o *Order) Cancel() error {
	if o.Status != OrderIsPending {
		return ErrOrderCannotBeCancelled
	}

	o.Status = OrderIsCancelled

	o.AddEvent(&OrderCanceled{
		Order: o,
	})

	return nil
}

func (o *Order) isPending() bool {
	return o.Status == OrderIsPending
}

func (o *Order) Ready() error {
	if !o.isPending() {
		return ErrOrderCannotBeReady
	}

	o.Status = OrderIsReady

	o.AddEvent(&OrderReadied{
		Order: o,
	})

	return nil
}

func (o *Order) isReady() bool {
	return o.Status == OrderIsReady
}

func (o *Order) Complete(invoiceID string) error {
	if invoiceID == "" {
		return ErrInvoiceIDCannotBeBlank
	}

	if !o.isReady() {
		return ErrOrderCannotBeCompleted
	}

	o.InvoiceID = invoiceID
	o.Status = OrderIsCompleted

	o.AddEvent(&OrderCompleted{
		Order: o,
	})

	return nil
}

func (o Order) GetTotal() float64 {
	var total float64
	for _, item := range o.Items {
		total += item.Price * float64(item.Quantity)
	}

	return total
}
