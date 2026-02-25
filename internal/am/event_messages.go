package am

import (
	"context"
	"mall/internal/ddd"
	"mall/internal/registry"
	"time"

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type EventMessage interface {
	MessageBase
	ddd.Event
}

type IncomingEventMessage interface {
	IncomingMessageBase
	ddd.Event
}

type eventMessage struct {
	id          string
	name        string
	payload     ddd.EventPayload
	occurred_at time.Time
	msg         IncomingMessageBase
}

var _ EventMessage = (*eventMessage)(nil)

func (e eventMessage) ID() string                { return e.id }
func (e eventMessage) EventName() string         { return e.name }
func (e eventMessage) Payload() ddd.EventPayload { return e.payload }
func (e eventMessage) Metadata() ddd.Metadata    { return e.msg.Metadata() }
func (e eventMessage) OccurredAt() time.Time     { return e.occurred_at }
func (e eventMessage) Subject() string           { return e.msg.Subject() }
func (e eventMessage) MessageName() string       { return e.msg.MessageName() }
func (e eventMessage) SentAt() time.Time         { return e.msg.SentAt() }
func (e eventMessage) ReceivedAt() time.Time     { return e.msg.ReceivedAt() }
func (e eventMessage) Ack() error                { return e.msg.Ack() }
func (e eventMessage) NAck() error               { return e.msg.NAck() }
func (e eventMessage) Extend() error             { return e.msg.Extend() }
func (e eventMessage) Kill() error               { return e.msg.Kill() }

type EventPublisher interface {
	Publish(ctx context.Context, topicName string, event ddd.Event) error
}

type eventPublisher struct {
	reg       registry.Registry
	publisher MessagePublisher
}

var _ EventPublisher = (*eventPublisher)(nil)

func NewEventPublisher(reg registry.Registry, msgPublisher MessagePublisher, mws ...MessagePublisherMiddleware) EventPublisher {
	return eventPublisher{
		reg:       reg,
		publisher: MessagePublisherWithMiddleware(msgPublisher, mws...),
	}
}

func (s eventPublisher) Publish(ctx context.Context, topicName string, event ddd.Event) error {
	payload, err := s.reg.Serialize(event.EventName(), event.Payload())
	if err != nil {
		return err
	}

	data, err := proto.Marshal(&EventMessageData{
		Payload:    payload,
		OccurredAt: timestamppb.New(event.OccurredAt()),
	})
	if err != nil {
		return err
	}

	return s.publisher.Publish(ctx, topicName, message{
		id:       event.ID(),
		name:     event.EventName(),
		subject:  topicName,
		data:     data,
		metadata: event.Metadata(),
		sentAt:   time.Now(),
	})
}

type eventMsgHandler struct {
	reg     registry.Registry
	handler ddd.EventHandler[ddd.Event]
}

var _ MessageHandler = (*eventMsgHandler)(nil)

func NewEventHandler(reg registry.Registry, handler ddd.EventHandler[ddd.Event], mws ...MessageHandlerMiddleware) MessageHandler {
	return MessageHandlerWithMiddleware(eventMsgHandler{
		reg:     reg,
		handler: handler,
	}, mws...)
}

func (h eventMsgHandler) HandleMessage(ctx context.Context, msg IncomingMessage) error {
	var eventData EventMessageData

	err := proto.Unmarshal(msg.Data(), &eventData)
	if err != nil {
		return err
	}

	eventName := msg.MessageName()

	payload, err := h.reg.Deserialize(eventName, eventData.GetPayload())
	if err != nil {
		return err
	}

	eventMsg := eventMessage{
		id:          msg.ID(),
		name:        eventName,
		payload:     payload,
		occurred_at: eventData.OccurredAt.AsTime(),
		msg:         msg,
	}

	return h.handler.HandleEvent(ctx, eventMsg)
}
