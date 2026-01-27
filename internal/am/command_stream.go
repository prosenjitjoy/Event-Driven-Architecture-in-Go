package am

import (
	"context"
	"mall/internal/ddd"
	"mall/internal/registry"
	"strings"

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/structpb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type CommandPublisher = MessagePublisher[ddd.Command]

type CommandSubscriber interface {
	Subscribe(topicName string, handler CommandMessageHandler, options ...SubscriberOption) error
}

type CommandStream interface {
	MessagePublisher[ddd.Command]
	CommandSubscriber
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
		id:   command.ID(),
		name: command.CommandName(),
		data: data,
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

		destination := commandMsg.Metadata().Get(CommandReplyChannelHeader).(string)

		reply, err := handler.HandleMessage(ctx, commandMsg)
		if err != nil {
			return s.publishReply(ctx, destination, s.failure(reply, commandMsg))
		}

		return s.publishReply(ctx, destination, s.success(reply, commandMsg))
	})

	return s.stream.Subscribe(topicName, fn, options...)
}

func (s commandStream) publishReply(ctx context.Context, destination string, reply ddd.Reply) error {
	metadata, err := structpb.NewStruct(reply.Metadata())
	if err != nil {
		return err
	}

	var payload []byte

	if reply.ReplyName() != SuccessReply && reply.ReplyName() != FailureReply {
		payload, err = s.reg.Serialize(reply.ReplyName(), reply.Payload())
		if err != nil {
			return err
		}
	}

	data, err := proto.Marshal(&ReplyMessageData{
		Payload:   payload,
		OccuredAt: timestamppb.New(reply.OccurredAt()),
		Metadata:  metadata,
	})
	if err != nil {
		return err
	}

	return s.stream.Publish(ctx, destination, rawMessage{
		id:   reply.ID(),
		name: reply.ReplyName(),
		data: data,
	})
}

func (s commandStream) failure(reply ddd.Reply, cmd ddd.Command) ddd.Reply {
	if reply == nil {
		reply = ddd.NewReply(FailureReply, nil)
	}

	reply.Metadata().Set(ReplyOutcomeHeader, OutcomeFailure)

	return s.applyCorrelationHeaders(reply, cmd)
}

func (s commandStream) success(reply ddd.Reply, cmd ddd.Command) ddd.Reply {
	if reply == nil {
		reply = ddd.NewReply(SuccessReply, nil)
	}

	reply.Metadata().Set(ReplyOutcomeHeader, OutcomeSuccess)

	return s.applyCorrelationHeaders(reply, cmd)
}

func (s commandStream) applyCorrelationHeaders(reply ddd.Reply, cmd ddd.Command) ddd.Reply {
	for key, value := range cmd.Metadata() {
		if key == CommandNameHeader {
			continue
		}

		if strings.HasPrefix(key, CommandHeaderPrefix) {
			hdr := ReplyHeaderPrefix + key[len(CommandHeaderPrefix):]

			reply.Metadata().Set(hdr, value)
		}
	}

	return reply
}
