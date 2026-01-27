package ddd

import (
	"time"

	"github.com/google/uuid"
)

type CommandOption interface {
	configureCommand(*command)
}

type CommandPayload any

type Command interface {
	IDer
	CommandName() string
	Payload() CommandPayload
	Metadata() Metadata
	OccurredAt() time.Time
}

type command struct {
	Entity
	payload    CommandPayload
	metadata   Metadata
	occurredAt time.Time
}

var _ Command = (*command)(nil)

func NewCommand(name string, payload CommandPayload, options ...CommandOption) Command {
	return newCommand(name, payload, options...)
}

func newCommand(name string, payload CommandPayload, options ...CommandOption) command {
	cmd := command{
		Entity:     NewEntity(uuid.New().String(), name),
		payload:    payload,
		metadata:   make(Metadata),
		occurredAt: time.Now(),
	}

	for _, option := range options {
		option.configureCommand(&cmd)
	}

	return cmd
}

func (c command) CommandName() string     { return c.name }
func (c command) Payload() CommandPayload { return c.payload }
func (c command) Metadata() Metadata      { return c.metadata }
func (c command) OccurredAt() time.Time   { return c.occurredAt }
