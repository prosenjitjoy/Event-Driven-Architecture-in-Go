package tm

import (
	"context"
	"errors"
	"mall/internal/am"
)

type OutboxStore interface {
	Save(ctx context.Context, msg am.Message) error
	FindUnpublished(ctx context.Context, limit int) ([]am.Message, error)
	MarkPublished(ctx context.Context, ids ...string) error
}

func OutboxPublisher(store OutboxStore) am.MessagePublisherMiddleware {
	return func(next am.MessagePublisher) am.MessagePublisher {
		return am.MessagePublisherFunc(func(ctx context.Context, topicName string, msg am.Message) error {
			var errDup ErrDuplicateMessage

			err := store.Save(ctx, msg)
			if errors.As(err, &errDup) {
				return nil
			}

			return err
		})
	}
}
