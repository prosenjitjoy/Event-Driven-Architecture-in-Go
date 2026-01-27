package am

import (
	"context"
	"mall/internal/ddd"
	"mall/internal/registry"

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/structpb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type ReplyPublisher = MessagePublisher[ddd.Reply]
type ReplySubscriber = MessageSubscriber[IncomingReplyMessage]
type ReplyStream = MessageStream[ddd.Reply, IncomingReplyMessage]

type replyStream struct {
	reg    registry.Registry
	stream RawMessageStream
}

var _ ReplyStream = (*replyStream)(nil)

func NewReplyStream(reg registry.Registry, stream RawMessageStream) ReplyStream {
	return &replyStream{
		reg:    reg,
		stream: stream,
	}
}

func (s replyStream) Publish(ctx context.Context, topicName string, reply ddd.Reply) error {
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

	return s.stream.Publish(ctx, topicName, rawMessage{
		id:   reply.ID(),
		name: reply.ReplyName(),
		data: data,
	})
}

func (s replyStream) Subscribe(topicName string, handler MessageHandler[IncomingReplyMessage], options ...SubscriberOption) error {
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

		var replyData ReplyMessageData

		err := proto.Unmarshal(msg.Data(), &replyData)
		if err != nil {
			return err
		}

		replyName := msg.MessageName()

		var payload any

		if replyName != SuccessReply && replyName != FailureReply {
			payload, err = s.reg.Deserialize(replyName, replyData.GetPayload())
			if err != nil {
				return err
			}
		}

		replyMsg := replyMessage{
			id:         msg.ID(),
			name:       replyName,
			payload:    payload,
			metadata:   replyData.Metadata.AsMap(),
			occurredAt: replyData.OccuredAt.AsTime(),
			msg:        msg,
		}

		return handler.HandleMessage(ctx, replyMsg)
	})

	return s.stream.Subscribe(topicName, fn, options...)
}
