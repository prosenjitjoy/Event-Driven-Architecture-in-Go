package logging

import (
	"context"
	"fmt"
	"log/slog"
	"mall/internal/ddd"
	"mall/internal/sec"
)

type ReplyHandlers[T any] struct {
	sec.Orchestrator[T]
	label  string
	logger *slog.Logger
}

var _ sec.Orchestrator[any] = (*ReplyHandlers[any])(nil)

func LogReplyHandlerAccess[T any](orchestrator sec.Orchestrator[T], label string, logger *slog.Logger) sec.Orchestrator[T] {
	return ReplyHandlers[T]{
		Orchestrator: orchestrator,
		label:        label,
		logger:       logger,
	}
}

func (h ReplyHandlers[T]) HandleReply(ctx context.Context, reply ddd.Reply) (err error) {
	h.logger.Info(fmt.Sprintf("--> COSEC.%s.On(%s)", h.label, reply.ReplyName()))
	defer func() {
		if err != nil {
			h.logger.Error(err.Error())
		}
		h.logger.Info(fmt.Sprintf("<-- COSEC.%s.On(%s)", h.label, reply.ReplyName()))
	}()

	return h.Orchestrator.HandleReply(ctx, reply)
}
