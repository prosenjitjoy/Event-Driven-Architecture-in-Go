package am

import (
	"context"
	"mall/internal/ddd"
	"mall/internal/registry"

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/structpb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type CommandPublisher = MessagePublisher[ddd.Command]
type CommandSubscriber = MessageSubscriber[IncomingCommandMessage]

type CommandStream interface {
	MessagePublisher[ddd.Command]
	MessageSubscriber[IncomingCommandMessage]
}

type commandStream struct {
	reg    registry.Registry
	stream RawMessageStream
}

var _ CommandStream = (*commandStream)(nil)

func NewCommandStream(reg registry.Registry, stream RawMessageStream) CommandStream {
	return &commandStream{
		reg:    reg,
		stream: stream,
	}
}

func (s commandStream) Publish(ctx context.Context, topicName string, command ddd.Command) error {
	metadata, err := structpb.NewStruct(command.Metadata())
	if err != nil {
		return err
	}

	payload, err := s.reg.Serialize(command.CommandName(), command.Payload())
	if err != nil {
		return err
	}

	data, err := proto.Marshal(&CommandMessageData{
		Payload:    payload,
		OccurredAt: timestamppb.New(command.OccurredAt()),
		Metadata:   metadata,
	})
	if err != nil {
		return err
	}

	return s.stream.Publish(ctx, topicName, rawMessage{
		id:      command.ID(),
		name:    command.CommandName(),
		subject: topicName,
		data:    data,
	})
}

func (s commandStream) Subscribe(topicName string, handler CommandMessageHandler, options ...SubscriberOption) error {
	cfg := NewSubscriberConfig(options)

	var filters map[string]struct{}

	if len(cfg.MessageFilters()) > 0 {
		filters = make(map[string]struct{})

		for _, key := range cfg.MessageFilters() {
			filters[key] = struct{}{}
		}
	}

	fn := MessageHandlerFunc[IncomingRawMessage](func(ctx context.Context, msg IncomingRawMessage) error {
		if filters != nil {
			if _, exists := filters[msg.MessageName()]; !exists {
				return nil
			}
		}

		var commandData CommandMessageData

		err := proto.Unmarshal(msg.Data(), &commandData)
		if err != nil {
			return err
		}

		commandName := msg.MessageName()

		payload, err := s.reg.Deserialize(commandName, commandData.GetPayload())
		if err != nil {
			return err
		}

		commandMsg := commandMessage{
			id:         msg.ID(),
			name:       commandName,
			payload:    payload,
			metadata:   commandData.GetMetadata().AsMap(),
			occurredAt: commandData.OccurredAt.AsTime(),
			msg:        msg,
		}

		return handler.HandleMessage(ctx, commandMsg)
	})

	return s.stream.Subscribe(topicName, fn, options...)
}
