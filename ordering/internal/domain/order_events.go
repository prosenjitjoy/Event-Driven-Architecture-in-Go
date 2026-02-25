package domain

type OrderCreated struct {
	CustomerID string
	PaymentID  string
	ShoppingID string
	Items      []Item
}

type OrderRejected struct{}

type OrderApproved struct {
	ShoppingID string
}

type OrderCanceled struct {
	CustomerID string
	PaymentID  string
}

type OrderReadied struct {
	CustomerID string
	PaymentID  string
	Total      float64
}

type OrderCompleted struct {
	CustomerID string
	InvoiceID  string
}
