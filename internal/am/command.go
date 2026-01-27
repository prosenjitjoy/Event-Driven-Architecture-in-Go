package am

import (
	"context"
	"mall/internal/ddd"
)

const (
	CommandHeaderPrefix       = "COMMAND_"
	CommandNameHeader         = CommandHeaderPrefix + "NAME"
	CommandReplyChannelHeader = CommandHeaderPrefix + "REPLY_CHANNEL"
)

type CommandMessageHandler interface {
	HandleMessage(ctx context.Context, msg IncomingCommandMessage) (ddd.Reply, error)
}

type CommandMessageHandlerFunc func(ctx context.Context, msg IncomingCommandMessage) (ddd.Reply, error)

func (f CommandMessageHandlerFunc) HandleMessage(ctx context.Context, cmd IncomingCommandMessage) (ddd.Reply, error) {
	return f(ctx, cmd)
}

type Command interface {
	ddd.Command
	Destination() string
}

type command struct {
	ddd.Command
	destination string
}

func NewCommand(name, destination string, payload ddd.CommandPayload, options ...ddd.CommandOption) Command {
	return command{
		Command:     ddd.NewCommand(name, payload, options...),
		destination: destination,
	}
}

func (c command) Destination() string {
	return c.destination
}
