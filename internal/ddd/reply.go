package ddd

import (
	"time"

	"github.com/google/uuid"
)

type ReplyOption interface {
	configureReply(*reply)
}

type ReplyPayload any

type Reply interface {
	IDer
	ReplyName() string
	Payload() ReplyPayload
	Metadata() Metadata
	OccurredAt() time.Time
}

type reply struct {
	Entity
	payload    ReplyPayload
	metadata   Metadata
	occurredAt time.Time
}

var _ Reply = (*reply)(nil)

func NewReply(name string, payload ReplyPayload, options ...ReplyOption) Reply {
	return newReply(name, payload, options...)
}

func newReply(name string, payload ReplyPayload, options ...ReplyOption) reply {
	rep := reply{
		Entity:     NewEntity(uuid.New().String(), name),
		payload:    payload,
		metadata:   make(Metadata),
		occurredAt: time.Now(),
	}

	for _, option := range options {
		option.configureReply(&rep)
	}

	return rep
}

func (r reply) ReplyName() string     { return r.name }
func (r reply) Payload() ReplyPayload { return r.payload }
func (r reply) Metadata() Metadata    { return r.metadata }
func (r reply) OccurredAt() time.Time { return r.occurredAt }
