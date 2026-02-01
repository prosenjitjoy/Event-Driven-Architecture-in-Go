package am

import (
	"context"
	"mall/internal/ddd"
	"mall/internal/registry"

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/structpb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type EventPublisher = MessagePublisher[ddd.Event]
type EventSubscriber = MessageSubscriber[IncomingEventMessage]
type EventStream = MessageStream[ddd.Event, IncomingEventMessage]

type eventStream struct {
	reg    registry.Registry
	stream RawMessageStream
}

var _ EventStream = (*eventStream)(nil)

func NewEventStream(reg registry.Registry, stream RawMessageStream) EventStream {
	return &eventStream{
		reg:    reg,
		stream: stream,
	}
}

func (s eventStream) Publish(ctx context.Context, topicName string, event ddd.Event) error {
	metadata, err := structpb.NewStruct(event.Metadata())
	if err != nil {
		return err
	}

	payload, err := s.reg.Serialize(event.EventName(), event.Payload())
	if err != nil {
		return err
	}

	data, err := proto.Marshal(&EventMessageData{
		Payload:    payload,
		OccurredAt: timestamppb.New(event.OccurredAt()),
		Metadata:   metadata,
	})
	if err != nil {
		return err
	}

	return s.stream.Publish(ctx, topicName, rawMessage{
		id:      event.ID(),
		name:    event.EventName(),
		subject: topicName,
		data:    data,
	})
}

func (s eventStream) Subscribe(topicName string, handler MessageHandler[IncomingEventMessage], options ...SubscriberOption) error {
	cfg := NewSubscriberConfig(options)

	var filters map[string]struct{}

	if len(cfg.MessageFilters()) > 0 {
		filters = make(map[string]struct{})
		for _, key := range cfg.MessageFilters() {
			filters[key] = struct{}{}
		}
	}

	fn := MessageHandlerFunc[IncomingRawMessage](func(ctx context.Context, msg IncomingRawMessage) error {
		var eventData EventMessageData

		if filters != nil {
			if _, exists := filters[msg.MessageName()]; !exists {
				return nil
			}
		}

		err := proto.Unmarshal(msg.Data(), &eventData)
		if err != nil {
			return nil
		}

		eventName := msg.MessageName()

		payload, err := s.reg.Deserialize(eventName, eventData.GetPayload())
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

		return handler.HandleMessage(ctx, eventMsg)
	})

	return s.stream.Subscribe(topicName, fn, options...)
}
