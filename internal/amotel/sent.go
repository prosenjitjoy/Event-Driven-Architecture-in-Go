package amotel

import (
	"context"
	"fmt"
	"mall/internal/am"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
)

func SentMessagesCounter(serviceName string) am.MessagePublisherMiddleware {
	counter, err := meter.Int64Counter(
		"sent_messages_count",
		metric.WithDescription(fmt.Sprintf("The total number of messages sent by %s", serviceName)),
	)
	if err != nil {
		panic(err)
	}

	return func(next am.MessagePublisher) am.MessagePublisher {
		return am.MessagePublisherFunc(func(ctx context.Context, topicName string, msg am.Message) (err error) {
			defer func() {
				msgName := msg.MessageName()

				// Increment the "all" bucket
				counter.Add(ctx, 1, metric.WithAttributes(
					attribute.String("message", "all"),
				))

				// Increment the specific message bucket
				counter.Add(ctx, 1, metric.WithAttributes(
					attribute.String("message", msgName),
				))
			}()

			return next.Publish(ctx, topicName, msg)
		})
	}
}
