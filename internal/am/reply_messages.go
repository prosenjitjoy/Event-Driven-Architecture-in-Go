package am

import (
	"context"
	"mall/internal/ddd"
	"mall/internal/registry"
	"time"

	"google.golang.org/protobuf/proto"
)

const (
	ReplyHeaderPrefix  = "REPLY_"
	ReplyNameHeader    = ReplyHeaderPrefix + "NAME"
	ReplyOutcomeHeader = ReplyHeaderPrefix + "OUTCOME"

	FailureReply = "am.Failure"
	SuccessReply = "am.Success"

	OutcomeSuccess = "SUCCESS"
	OutcomeFailure = "FAILURE"
)

type ReplyMessage interface {
	Message
	ddd.Reply
}

type IncomingReplyMessage interface {
	IncomingMessage
	ddd.Reply
}

type replyMessage struct {
	id         string
	name       string
	payload    ddd.ReplyPayload
	metadata   ddd.Metadata
	occurredAt time.Time
	msg        IncomingMessage
}

var _ ReplyMessage = (*replyMessage)(nil)

func (r replyMessage) ID() string                { return r.id }
func (r replyMessage) ReplyName() string         { return r.name }
func (r replyMessage) Payload() ddd.ReplyPayload { return r.payload }
func (r replyMessage) Metadata() ddd.Metadata    { return r.metadata }
func (r replyMessage) OccurredAt() time.Time     { return r.occurredAt }
func (r replyMessage) Subject() string           { return r.msg.Subject() }
func (r replyMessage) MessageName() string       { return r.msg.MessageName() }
func (r replyMessage) Ack() error                { return r.msg.Ack() }
func (r replyMessage) NAck() error               { return r.msg.NAck() }
func (r replyMessage) Extend() error             { return r.msg.Extend() }
func (r replyMessage) Kill() error               { return r.msg.Kill() }

type replyMsgHandler struct {
	reg     registry.Registry
	handler ddd.ReplyHandler[ddd.Reply]
}

func NewReplyMessageHandler(reg registry.Registry, handler ddd.ReplyHandler[ddd.Reply]) RawMessageHandler {
	return replyMsgHandler{
		reg:     reg,
		handler: handler,
	}
}

func (h replyMsgHandler) HandleMessage(ctx context.Context, msg IncomingRawMessage) error {
	var replyData ReplyMessageData

	err := proto.Unmarshal(msg.Data(), &replyData)
	if err != nil {
		return err
	}

	replyName := msg.MessageName()

	var payload any

	if replyName != SuccessReply && replyName != FailureReply {
		payload, err = h.reg.Deserialize(replyName, replyData.GetPayload())
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

	return h.handler.HandleReply(ctx, replyMsg)
}
