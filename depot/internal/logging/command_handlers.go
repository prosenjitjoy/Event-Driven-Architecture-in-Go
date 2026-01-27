package logging

import (
	"context"
	"fmt"
	"log/slog"
	"mall/internal/ddd"
)

type CommandHandlers[T ddd.Command] struct {
	ddd.CommandHandler[T]
	label  string
	logger *slog.Logger
}

var _ ddd.CommandHandler[ddd.Command] = (*CommandHandlers[ddd.Command])(nil)

func LogCommandHandlerAccess[T ddd.Command](handlers ddd.CommandHandler[T], label string, logger *slog.Logger) ddd.CommandHandler[T] {
	return CommandHandlers[T]{
		CommandHandler: handlers,
		label:          label,
		logger:         logger,
	}
}

func (h CommandHandlers[T]) HandleCommand(ctx context.Context, command T) (reply ddd.Reply, err error) {
	h.logger.Info(fmt.Sprintf("--> Depot.%s.On(%s)", h.label, command.CommandName()))
	defer func() {
		if err != nil {
			h.logger.Error(err.Error())
		}
		h.logger.Info(fmt.Sprintf("<-- Depot.%s.On(%s)", h.label, command.CommandName()))
	}()

	return h.CommandHandler.HandleCommand(ctx, command)
}
