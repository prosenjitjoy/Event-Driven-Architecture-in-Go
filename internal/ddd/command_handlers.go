package ddd

import "context"

type CommandHandler[T Command] interface {
	HandleCommand(ctx context.Context, cmd T) (Reply, error)
}

type CommandHandlerFunc[T Command] func(ctx context.Context, cmd T) (Reply, error)

func (f CommandHandlerFunc[T]) HandleCommand(ctx context.Context, cmd T) (Reply, error) {
	return f(ctx, cmd)
}
