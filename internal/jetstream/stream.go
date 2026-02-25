package jetstream

import (
	"context"
	"fmt"
	"log/slog"
	"mall/internal/am"
	"sync"
	"time"

	"github.com/nats-io/nats.go"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/structpb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

const maxRetries = 5

type Stream struct {
	streamName string
	js         nats.JetStreamContext
	subs       []*nats.Subscription
	logger     *slog.Logger
	mu         sync.Mutex
}

var _ am.MessageStream = (*Stream)(nil)

func NewStream(streamName string, js nats.JetStreamContext, logger *slog.Logger) *Stream {
	return &Stream{
		streamName: streamName,
		js:         js,
		logger:     logger,
	}
}

func (s *Stream) Publish(ctx context.Context, topicName string, rawMsg am.Message) error {
	metadata, err := structpb.NewStruct(rawMsg.Metadata())
	if err != nil {
		return err
	}

	data, err := proto.Marshal(&StreamMessage{
		Id:       rawMsg.ID(),
		Name:     rawMsg.MessageName(),
		Data:     rawMsg.Data(),
		Metadata: metadata,
		SentAt:   timestamppb.New(rawMsg.SentAt()),
	})
	if err != nil {
		return err
	}

	p, err := s.js.PublishMsgAsync(&nats.Msg{
		Subject: topicName,
		Data:    data,
	}, nats.MsgId(rawMsg.ID()))
	if err != nil {
		return err
	}

	// retry a handful of times to publish the messages
	go func(ctx context.Context, future nats.PubAckFuture, tries int) {
		var err error

		for {
			select {
			case <-future.Ok(): // publish acknowledged
				return
			case <-future.Err(): // error ignored: try again
				// TODO: add variable delay between tries
				tries = tries - 1
				if tries <= 0 {
					s.logger.ErrorContext(ctx, fmt.Sprintf("unable to publish message after %d tries", maxRetries))
					return
				}
				future, err = s.js.PublishMsgAsync(future.Msg())
				if err != nil {
					s.logger.ErrorContext(ctx, fmt.Sprintf("failed to publish a message: %s", err.Error()))
					return
				}
			}
		}
	}(ctx, p, maxRetries)

	return nil
}

func (s *Stream) Subscribe(topicName string, handler am.MessageHandler, options ...am.SubscriberOption) (am.Subscription, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	subCfg := am.NewSubscriberConfig(options)

	opts := []nats.SubOpt{
		nats.MaxDeliver(subCfg.MaxRedeliver()),
	}

	cfg := &nats.ConsumerConfig{
		MaxDeliver:     subCfg.MaxRedeliver(),
		DeliverSubject: topicName,
		FilterSubject:  topicName,
	}

	if groupName := subCfg.GroupName(); groupName != "" {
		cfg.DeliverSubject = groupName
		cfg.DeliverGroup = groupName
		cfg.Durable = groupName

		opts = append(opts, nats.Bind(s.streamName, groupName), nats.Durable(groupName))
	}

	if ackType := subCfg.AckType(); ackType != am.AckTypeAuto {
		ackWait := subCfg.AckWait()

		cfg.AckPolicy = nats.AckExplicitPolicy
		cfg.AckWait = ackWait

		opts = append(opts, nats.AckExplicit(), nats.AckWait(ackWait))
	} else {
		cfg.AckPolicy = nats.AckNonePolicy
		opts = append(opts, nats.AckNone())
	}

	_, err := s.js.AddConsumer(s.streamName, cfg)
	if err != nil {
		return nil, err
	}

	var sub *nats.Subscription

	if groupName := subCfg.GroupName(); groupName == "" {
		sub, err = s.js.Subscribe(topicName, s.handleMsg(subCfg, handler), opts...)
	} else {
		sub, err = s.js.QueueSubscribe(topicName, groupName, s.handleMsg(subCfg, handler), opts...)
	}

	s.subs = append(s.subs, sub)

	return subscription{sub}, nil
}

func (s *Stream) Unsubscribe() error {
	for _, sub := range s.subs {
		if !sub.IsValid() {
			continue
		}

		err := sub.Drain()
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *Stream) handleMsg(cfg am.SubscriberConfig, handler am.MessageHandler) func(*nats.Msg) {
	var filters map[string]struct{}

	if len(cfg.MessageFilters()) > 0 {
		filters = make(map[string]struct{})

		for _, key := range cfg.MessageFilters() {
			filters[key] = struct{}{}
		}
	}

	return func(natsMsg *nats.Msg) {
		m := &StreamMessage{}

		err := proto.Unmarshal(natsMsg.Data, m)
		if err != nil {
			s.logger.WarnContext(context.TODO(), "failed to unmarshal the *nats.Msg", "error", err.Error())

			return
		}

		if filters != nil {
			if _, exists := filters[m.GetName()]; !exists {
				err := natsMsg.Ack()
				if err != nil {
					s.logger.WarnContext(context.TODO(), "failed to Ack a filtered message")
				}

				return
			}
		}

		msg := &rawMessage{
			id:         m.GetId(),
			name:       m.GetName(),
			subject:    natsMsg.Subject,
			data:       m.GetData(),
			metadata:   m.GetMetadata().AsMap(),
			sentAt:     m.SentAt.AsTime(),
			receivedAt: time.Now(),
			acked:      false,
			ackFn:      func() error { return natsMsg.Ack() },
			nackFn:     func() error { return natsMsg.Nak() },
			extendFn:   func() error { return natsMsg.InProgress() },
			killFn:     func() error { return natsMsg.Term() },
		}

		wCtx, cancel := context.WithTimeout(context.Background(), cfg.AckWait())
		defer cancel()

		errc := make(chan error)
		go func() {
			errc <- handler.HandleMessage(wCtx, msg)
		}()

		if cfg.AckType() == am.AckTypeAuto {
			err := msg.Ack()
			if err != nil {
				s.logger.WarnContext(context.TODO(), "failed to auto-Ack a message", "error", err.Error())
			}
		}

		select {
		case err = <-errc:
			if err == nil {
				if ackErr := msg.Ack(); ackErr != nil {
					s.logger.WarnContext(context.TODO(), "failed to Ack a message", "error", ackErr.Error())
				}
				return
			}
			if nakErr := msg.NAck(); nakErr != nil {
				s.logger.WarnContext(context.TODO(), "failed to Nack a message", "error", nakErr.Error())
			}

			s.logger.WarnContext(context.TODO(), "error while handling message", "error", err.Error())
		case <-wCtx.Done():
			return
		}
	}
}
