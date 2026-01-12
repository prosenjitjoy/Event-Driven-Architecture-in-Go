package ddd

import (
	"context"
)

type Event interface {
	EventName() string
}

type EventHandler func(ctx context.Context, event Event) error
