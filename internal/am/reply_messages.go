package am

import (
	"mall/internal/ddd"
	"time"
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
func (r replyMessage) MessageName() string       { return r.msg.MessageName() }
func (r replyMessage) Ack() error                { return r.msg.Ack() }
func (r replyMessage) NAck() error               { return r.msg.NAck() }
func (r replyMessage) Extend() error             { return r.msg.Extend() }
func (r replyMessage) Kill() error               { return r.msg.Kill() }
