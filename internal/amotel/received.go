package amotel

import (
	"context"
	"fmt"
	"mall/internal/am"
	"strconv"
	"time"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
)

func ReceivedMessagesCounter(serviceName string) am.MessageHandlerMiddleware {
	counter, err := meter.Int64Counter(
		"received_messages_count",
		metric.WithDescription(fmt.Sprintf("The total number of messages received by %s", serviceName)),
	)
	if err != nil {
		panic(err)
	}

	histogram, err := meter.Float64Histogram(
		"received_messages_latency_seconds",
		metric.WithDescription("Latency of received messages"),
	)
	if err != nil {
		panic(err)
	}

	return func(next am.MessageHandler) am.MessageHandler {
		return am.MessageHandlerFunc(func(ctx context.Context, msg am.IncomingMessage) (err error) {
			defer func(started time.Time) {
				elapsed := time.Since(started).Seconds()
				handled := strconv.FormatBool(err == nil)
				msgName := msg.MessageName()

				commonAttrs := []attribute.KeyValue{
					attribute.String("handled", handled),
				}

				counter.Add(ctx, 1, metric.WithAttributes(
					append(commonAttrs, attribute.String("message", "all"))...,
				))
				counter.Add(ctx, 1, metric.WithAttributes(
					append(commonAttrs, attribute.String("message", msgName))...,
				))

				histogram.Record(ctx, elapsed, metric.WithAttributes(
					append(commonAttrs, attribute.String("message", "all"))...,
				))
				histogram.Record(ctx, elapsed, metric.WithAttributes(
					append(commonAttrs, attribute.String("message", msgName))...,
				))
			}(time.Now())

			return next.HandleMessage(ctx, msg)
		})
	}
}
