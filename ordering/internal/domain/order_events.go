package domain

type OrderCreated struct {
	Order *Order
}

type OrderCanceled struct {
	Order *Order
}

type OrderReadied struct {
	Order *Order
}

type OrderCompleted struct {
	Order *Order
}

func (OrderCreated) EventName() string {
	return "ordering.OrderCreated"
}

func (OrderCanceled) EventName() string {
	return "ordering.OrderCanceled"
}

func (OrderReadied) EventName() string {
	return "ordering.OrderReadied"
}

func (OrderCompleted) EventName() string {
	return "ordering.OrderCompleted"
}
