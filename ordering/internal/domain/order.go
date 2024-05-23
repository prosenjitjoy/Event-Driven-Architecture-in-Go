package domain

import (
	"errors"
	"mall/internal/ddd"
)

var (
	ErrOrderHasNoItems         = errors.New("order has no items")
	ErrOrderCannotBeCancelled  = errors.New("order cannot be cancelled")
	ErrCustomerIDCannotBeBlank = errors.New("customer id cannot be blank")
	ErrPaymentIDCannotBeBlank  = errors.New("payment id cannot be blank")
	ErrOrderCannotBeReady      = errors.New("order cannot be ready")
	ErrOrderCannotBeCompleted  = errors.New("order cannot be completed")
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

	order.AddEvent(&OrderCreated{Order: order})

	return order, nil
}

func (o *Order) Cancel() error {
	if o.Status != OrderIsPending {
		return ErrOrderCannotBeCancelled
	}

	o.Status = OrderIsCanceled

	o.AddEvent(&OrderCanceled{Order: o})

	return nil
}

func (o *Order) Ready() error {
	if o.Status != OrderIsPending {
		return ErrOrderCannotBeReady
	}

	o.Status = OrderIsReady

	o.AddEvent(&OrderReadied{Order: o})

	return nil
}

func (o *Order) Complete(invoiceID string) error {
	// validate invoice exists

	if o.Status != OrderIsReady {
		return ErrOrderCannotBeCompleted
	}

	o.InvoiceID = invoiceID
	o.Status = OrderIsCompleted

	o.AddEvent(&OrderCompleted{Order: o})

	return nil
}

func (o Order) GetTotal() float64 {
	var total float64

	for _, item := range o.Items {
		total += item.ProductPrice * float64(item.Quantity)
	}

	return total
}
