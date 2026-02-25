package jetstream

import (
	"mall/internal/am"

	"github.com/nats-io/nats.go"
)

type subscription struct {
	sub *nats.Subscription
}

var _ am.Subscription = (*subscription)(nil)

func (s subscription) Unsubscribe() error {
	if s.sub.IsValid() {
		return nil
	}

	return s.sub.Drain()
}
