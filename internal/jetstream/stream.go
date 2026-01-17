package jetstream

import (
	"context"
	"mall/internal/am"

	"github.com/nats-io/nats.go"
	"google.golang.org/protobuf/proto"
)

const maxRetries = 5

type Stream struct {
	streamName string
	js         nats.JetStreamContext
}

var _ am.MessageStream[am.RawMessage, am.RawMessage] = (*Stream)(nil)

func NewStream(streamName string, js nats.JetStreamContext) *Stream {
	return &Stream{
		streamName: streamName,
		js:         js,
	}
}

func (s *Stream) Publish(ctx context.Context, topicName string, rawMsg am.RawMessage) error {
	data, err := proto.Marshal(&StreamMessage{
		Id:   rawMsg.ID(),
		Name: rawMsg.MessageName(),
		Data: rawMsg.Data(),
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
	go func(future nats.PubAckFuture, tries int) {
		var err error

		for {
			select {
			case <-future.Ok(): // publish acknowledged
				return
			case <-future.Err(): // error: try again
				tries = tries - 1
				if tries <= 0 {
					return
				}
				future, err = s.js.PublishMsgAsync(future.Msg())
				if err != nil {
					return
				}
			}
		}
	}(p, maxRetries)

	return nil
}

func (s *Stream) Subscribe(topicName string, handler am.MessageHandler[am.RawMessage], options ...am.SubscriberOption) error {
	subCfg := am.NewSubscriberConfig(options)

	opts := []nats.SubOpt{
		nats.MaxDeliver(subCfg.MaxRedeliver()),
	}

	cfg := &nats.ConsumerConfig{
		MaxDeliver: subCfg.MaxRedeliver(),
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
		return err
	}

	if groupName := subCfg.GroupName(); groupName == "" {
		_, err := s.js.Subscribe(topicName, s.handleMsg(subCfg, handler), opts...)
		if err != nil {
			return err
		}
	} else {
		_, err := s.js.QueueSubscribe(topicName, groupName, s.handleMsg(subCfg, handler), opts...)
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *Stream) handleMsg(cfg am.SubscriberConfig, handler am.MessageHandler[am.RawMessage]) func(*nats.Msg) {
	return func(natsMsg *nats.Msg) {
		m := &StreamMessage{}
		err := proto.Unmarshal(natsMsg.Data, m)
		if err != nil {
			// NAck? ... Logging?
			return
		}

		msg := &rawMessage{
			id:       m.GetId(),
			name:     m.GetName(),
			data:     m.GetData(),
			acked:    false,
			ackFn:    func() error { return natsMsg.Ack() },
			nackFn:   func() error { return natsMsg.Nak() },
			extendFn: func() error { return natsMsg.InProgress() },
			killFn:   func() error { return natsMsg.Term() },
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
				// can be logged
			}
		}

		select {
		case err = <-errc:
			if err == nil {
				if ackErr := msg.Ack(); ackErr != nil {
					// can be logged
				}
				return
			}
			if nakErr := msg.NAck(); nakErr != nil {
				// can be logged
			}
		case <-wCtx.Done():
			return
		}
	}
}
