package domain

type OrderStatus string

const (
	OrderIsUnknown   OrderStatus = ""
	OrderIsPending   OrderStatus = "pending"
	OrderIsInProcess OrderStatus = "in-progress"
	OrderIsReady     OrderStatus = "ready"
	OrderIsCompleted OrderStatus = "completed"
	OrderIsCanceled  OrderStatus = "canceled"
)

func (s OrderStatus) String() string {
	switch s {
	case OrderIsPending, OrderIsInProcess, OrderIsReady, OrderIsCompleted, OrderIsCanceled:
		return string(s)
	default:
		return ""
	}
}

func ToOrderStatus(status string) OrderStatus {
	switch status {
	case OrderIsPending.String():
		return OrderIsPending
	case OrderIsInProcess.String():
		return OrderIsInProcess
	case OrderIsReady.String():
		return OrderIsReady
	case OrderIsCanceled.String():
		return OrderIsCanceled
	case OrderIsCompleted.String():
		return OrderIsCompleted
	default:
		return OrderIsUnknown
	}
}
