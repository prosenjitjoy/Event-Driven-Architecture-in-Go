package am

import (
	"context"
	"mall/internal/ddd"
	"mall/internal/registry"
	"strings"

	"google.golang.org/protobuf/proto"
)

type Command interface {
	ddd.Command
	Destination() string
}

type command struct {
	ddd.Command
	destination string
}

var _ Command = (*command)(nil)

func NewCommand(name, destination string, payload ddd.CommandPayload, options ...ddd.CommandOption) Command {
	return command{
		Command:     ddd.NewCommand(name, payload, options...),
		destination: destination,
	}
}

func (c command) Destination() string {
	return c.destination
}

type CommandMessageHandler = MessageHandler[IncomingCommandMessage]
type CommandMessageHandlerFunc func(ctx context.Context, msg IncomingCommandMessage) (ddd.Reply, error)

func (f CommandMessageHandlerFunc) HandleMessage(ctx context.Context, cmd IncomingCommandMessage) (ddd.Reply, error) {
	return f(ctx, cmd)
}

type commandMsgHandler struct {
	reg       registry.Registry
	publisher ReplyPublisher
	handler   ddd.CommandHandler[ddd.Command]
}

var _ RawMessageHandler = (*commandMsgHandler)(nil)

func NewCommandMessageHandler(reg registry.Registry, publisher ReplyPublisher, handler ddd.CommandHandler[ddd.Command]) RawMessageHandler {
	return commandMsgHandler{
		reg:       reg,
		publisher: publisher,
		handler:   handler,
	}
}

func (h commandMsgHandler) HandleMessage(ctx context.Context, msg IncomingRawMessage) error {
	var commandData CommandMessageData

	err := proto.Unmarshal(msg.Data(), &commandData)
	if err != nil {
		return err
	}

	commandName := msg.MessageName()

	payload, err := h.reg.Deserialize(commandName, commandData.GetPayload())
	if err != nil {
		return err
	}

	commandMsg := commandMessage{
		id:         msg.ID(),
		name:       commandName,
		payload:    payload,
		metadata:   commandData.Metadata.AsMap(),
		occurredAt: commandData.OccurredAt.AsTime(),
		msg:        msg,
	}

	destination := commandMsg.Metadata().Get(CommandReplyChannelHeader).(string)

	reply, err := h.handler.HandleCommand(ctx, commandMsg)
	if err != nil {
		return h.publisher.Publish(ctx, destination, h.failure(reply, commandMsg))
	}

	err = h.publisher.Publish(ctx, destination, h.success(reply, commandMsg))
	if err != nil {
		return err
	}

	return nil
}

func (s commandMsgHandler) failure(reply ddd.Reply, cmd ddd.Command) ddd.Reply {
	if reply == nil {
		reply = ddd.NewReply(FailureReply, nil)
	}

	reply.Metadata().Set(ReplyOutcomeHeader, OutcomeFailure)

	return s.applyCorrelationHeaders(reply, cmd)
}

func (s commandMsgHandler) success(reply ddd.Reply, cmd ddd.Command) ddd.Reply {
	if reply == nil {
		reply = ddd.NewReply(SuccessReply, nil)
	}

	reply.Metadata().Set(ReplyOutcomeHeader, OutcomeSuccess)

	return s.applyCorrelationHeaders(reply, cmd)
}

func (s commandMsgHandler) applyCorrelationHeaders(reply ddd.Reply, cmd ddd.Command) ddd.Reply {
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
