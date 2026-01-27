package am

import (
	"mall/internal/ddd"
	"time"
)

type CommandMessage interface {
	Message
	ddd.Command
}

type IncomingCommandMessage interface {
	IncomingMessage
	ddd.Command
}

type commandMessage struct {
	id         string
	name       string
	payload    ddd.CommandPayload
	metadata   ddd.Metadata
	occurredAt time.Time
	msg        IncomingMessage
}

var _ CommandMessage = (*commandMessage)(nil)

func (c commandMessage) ID() string                  { return c.id }
func (c commandMessage) CommandName() string         { return c.name }
func (c commandMessage) Payload() ddd.CommandPayload { return c.payload }
func (c commandMessage) Metadata() ddd.Metadata      { return c.metadata }
func (c commandMessage) OccurredAt() time.Time       { return c.occurredAt }
func (c commandMessage) MessageName() string         { return c.MessageName() }
func (c commandMessage) Ack() error                  { return c.msg.Ack() }
func (c commandMessage) NAck() error                 { return c.msg.NAck() }
func (c commandMessage) Extend() error               { return c.msg.Extend() }
func (c commandMessage) Kill() error                 { return c.msg.Kill() }
