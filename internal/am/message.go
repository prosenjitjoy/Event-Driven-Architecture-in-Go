package am

import (
	"context"
	"mall/internal/ddd"
	"time"
)

type MessageBase interface {
	ddd.IDer
	Subject() string
	MessageName() string
	Metadata() ddd.Metadata
	SentAt() time.Time
}

type IncomingMessageBase interface {
	MessageBase
	ReceivedAt() time.Time
	Ack() error
	NAck() error
	Extend() error
	Kill() error
}

type Message interface {
	MessageBase
	Data() []byte
}

type IncomingMessage interface {
	IncomingMessageBase
	Data() []byte
}

type MessageHandler interface {
	HandleMessage(ctx context.Context, msg IncomingMessage) error
}

type MessageHandlerFunc func(ctx context.Context, msg IncomingMessage) error

func (f MessageHandlerFunc) HandleMessage(ctx context.Context, cmd IncomingMessage) error {
	return f(ctx, cmd)
}

type MessageSubscriber interface {
	Subscribe(topicName string, handler MessageHandler, options ...SubscriberOption) (Subscription, error)
	Unsubscribe() error
}

type MessagePublisher interface {
	Publish(ctx context.Context, topicName string, msg Message) error
}

type MessagePublisherFunc func(ctx context.Context, topicName string, msg Message) error

func (f MessagePublisherFunc) Publish(ctx context.Context, topicName string, msg Message) error {
	return f(ctx, topicName, msg)
}

type MessageStream interface {
	MessageSubscriber
	MessagePublisher
}

type MessageStreamMiddleware = func(next MessageStream) MessageStream
type MessagePublisherMiddleware = func(next MessagePublisher) MessagePublisher
type MessageHandlerMiddleware = func(next MessageHandler) MessageHandler

type message struct {
	id       string
	name     string
	subject  string
	data     []byte
	metadata ddd.Metadata
	sentAt   time.Time
}

type messagePublisher struct {
	publisher MessagePublisher
}

type messageSubscriber struct {
	subscriber MessageSubscriber
	mws        []MessageHandlerMiddleware
}

var _ Message = (*message)(nil)

func (m message) ID() string             { return m.id }
func (m message) Subject() string        { return m.subject }
func (m message) MessageName() string    { return m.name }
func (m message) Data() []byte           { return m.data }
func (m message) Metadata() ddd.Metadata { return m.metadata }
func (m message) SentAt() time.Time      { return m.sentAt }

func NewMessagePublisher(publisher MessagePublisher, mws ...MessagePublisherMiddleware) MessagePublisher {
	return messagePublisher{
		publisher: MessagePublisherWithMiddleware(publisher, mws...),
	}
}

func (p messagePublisher) Publish(ctx context.Context, topicName string, msg Message) error {
	return p.publisher.Publish(ctx, topicName, msg)
}

func NewMessageSubscriber(subscriber MessageSubscriber, mws ...MessageHandlerMiddleware) MessageSubscriber {
	return messageSubscriber{
		subscriber: subscriber,
		mws:        mws,
	}
}

func (s messageSubscriber) Subscribe(topicName string, handler MessageHandler, option ...SubscriberOption) (Subscription, error) {
	return s.subscriber.Subscribe(topicName, MessageHandlerWithMiddleware(handler, s.mws...), option...)
}

func (s messageSubscriber) Unsubscribe() error {
	return s.subscriber.Unsubscribe()
}
