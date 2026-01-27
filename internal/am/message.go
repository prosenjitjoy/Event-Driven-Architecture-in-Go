package am

import (
	"context"
	"mall/internal/ddd"
)

type Message interface {
	ddd.IDer
	MessageName() string
}

type IncomingMessage interface {
	Message
	Ack() error
	NAck() error
	Extend() error
	Kill() error
}

type MessageHandler[I IncomingMessage] interface {
	HandleMessage(ctx context.Context, msg I) error
}

type MessageHandlerFunc[I IncomingMessage] func(ctx context.Context, msg I) error

func (f MessageHandlerFunc[I]) HandleMessage(ctx context.Context, msg I) error {
	return f(ctx, msg)
}

type MessagePublisher[O any] interface {
	Publish(ctx context.Context, topicName string, v O) error
}

type MessageSubscriber[I IncomingMessage] interface {
	Subscribe(topicName string, handler MessageHandler[I], options ...SubscriberOption) error
}

type MessageStream[O any, I IncomingMessage] interface {
	MessagePublisher[O]
	MessageSubscriber[I]
}
