package am

type RawMessageStream = MessageStream[RawMessage, IncomingRawMessage]
type RawMessageHandler = MessageHandler[IncomingRawMessage]
type RawMessagePublisher = MessagePublisher[RawMessage]
type RawMessageSubscriber = MessageSubscriber[IncomingRawMessage]

type RawMessage interface {
	Message
	Data() []byte
}

type IncomingRawMessage interface {
	IncomingMessage
	Data() []byte
}

type rawMessage struct {
	id   string
	name string
	data []byte
}

var _ RawMessage = (*rawMessage)(nil)

func (m rawMessage) ID() string          { return m.id }
func (m rawMessage) MessageName() string { return m.name }
func (m rawMessage) Data() []byte        { return m.data }
