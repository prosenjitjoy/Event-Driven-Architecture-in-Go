package am

import (
	"context"
	"mall/internal/ddd"
)

type Message interface {
	ddd.IDer
	MessageName() string
	Ack() error
	NAck() error
	Extend() error
	Kill() error
}

type MessageHandler[O Message] interface {
	HandleMessage(ctx context.Context, msg O) error
}

type MessageHandlerFunc[O Message] func(ctx context.Context, msg O) error

type MessagePublisher[I any] interface {
	Publish(ctx context.Context, topicName string, v I) error
}

type MessageSubscriber[O Message] interface {
	Subscribe(topicName string, handler MessageHandler[O], options ...SubscriberOption) error
}

type MessageStream[I any, O Message] interface {
	MessagePublisher[I]
	MessageSubscriber[O]
}

func (f MessageHandlerFunc[O]) HandleMessage(ctx context.Context, msg O) error {
	return f(ctx, msg)
}
