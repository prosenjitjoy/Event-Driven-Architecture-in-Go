package am

import "context"

type RawMessageStream = MessageStream[RawMessage, IncomingRawMessage]

type RawMessagePublisher = MessagePublisher[RawMessage]
type RawMessageSubscriber = MessageSubscriber[IncomingRawMessage]

type RawMessageHandler = MessageHandler[IncomingRawMessage]
type RawMessageHandlerFunc func(ctx context.Context, msg IncomingRawMessage) error

type RawMessageStreamMiddleware = func(stream RawMessageStream) RawMessageStream
type RawMessageHandlerMiddleware = func(handler RawMessageHandler) RawMessageHandler

type RawMessage interface {
	Message
	Data() []byte
}

type IncomingRawMessage interface {
	IncomingMessage
	Data() []byte
}

type rawMessage struct {
	id      string
	name    string
	subject string
	data    []byte
}

var _ RawMessage = (*rawMessage)(nil)

func (m rawMessage) ID() string          { return m.id }
func (m rawMessage) Subject() string     { return m.subject }
func (m rawMessage) MessageName() string { return m.name }
func (m rawMessage) Data() []byte        { return m.data }

func (f RawMessageHandlerFunc) HandleMessage(ctx context.Context, cmd IncomingRawMessage) error {
	return f(ctx, cmd)
}

func RawMessageStreamWithMiddleware(stream RawMessageStream, mws ...RawMessageStreamMiddleware) RawMessageStream {
	s := stream

	for i := len(mws) - 1; i >= 0; i-- {
		s = mws[i](s)
	}

	return s
}

func RawMessageHandlerWithMiddleware(handler RawMessageHandler, mws ...RawMessageHandlerMiddleware) RawMessageHandler {
	h := handler

	for i := len(mws) - 1; i >= 0; i-- {
		h = mws[i](h)
	}

	return h
}
