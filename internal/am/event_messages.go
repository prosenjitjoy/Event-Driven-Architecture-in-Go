package am

import (
	"context"
	"mall/internal/ddd"
	"mall/internal/registry"
	"time"

	"google.golang.org/protobuf/proto"
)

type EventMessage interface {
	Message
	ddd.Event
}

type IncomingEventMessage interface {
	IncomingMessage
	ddd.Event
}

type eventMessage struct {
	id          string
	name        string
	payload     ddd.EventPayload
	metadata    ddd.Metadata
	occurred_at time.Time
	msg         IncomingMessage
}

var _ EventMessage = (*eventMessage)(nil)

func (e eventMessage) ID() string                { return e.id }
func (e eventMessage) EventName() string         { return e.name }
func (e eventMessage) Payload() ddd.EventPayload { return e.payload }
func (e eventMessage) Metadata() ddd.Metadata    { return e.metadata }
func (e eventMessage) OccurredAt() time.Time     { return e.occurred_at }
func (e eventMessage) Subject() string           { return e.msg.Subject() }
func (e eventMessage) MessageName() string       { return e.msg.MessageName() }
func (e eventMessage) Ack() error                { return e.msg.Ack() }
func (e eventMessage) NAck() error               { return e.msg.NAck() }
func (e eventMessage) Extend() error             { return e.msg.Extend() }
func (e eventMessage) Kill() error               { return e.msg.Kill() }

type eventMsgHandler struct {
	reg     registry.Registry
	handler ddd.EventHandler[ddd.Event]
}

func NewEventMessageHandler(reg registry.Registry, handler ddd.EventHandler[ddd.Event]) RawMessageHandler {
	return eventMsgHandler{
		reg:     reg,
		handler: handler,
	}
}

func (h eventMsgHandler) HandleMessage(ctx context.Context, msg IncomingRawMessage) error {
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
		metadata:    eventData.GetMetadata().AsMap(),
		occurred_at: eventData.GetOccurredAt().AsTime(),
		msg:         msg,
	}

	return h.handler.HandleEvent(ctx, eventMsg)
}
