package waiter

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"golang.org/x/sync/errgroup"
)

type WaitFunc func(ctx context.Context) error
type CleanupFunc func()

type Waiter interface {
	AddWaitFunc(fns ...WaitFunc)
	AddCleanupFunc(fns ...CleanupFunc)
	Wait() error
	Context() context.Context
	CancelFunc() context.CancelFunc
}

type waiter struct {
	ctx          context.Context
	cancel       context.CancelFunc
	waitFuncs    []WaitFunc
	cleanupFuncs []CleanupFunc
}

type waiterCfg struct {
	parentCtx    context.Context
	catchSignals bool
}

func New(options ...WaiterOption) Waiter {
	cfg := &waiterCfg{
		parentCtx:    context.Background(),
		catchSignals: false,
	}

	for _, option := range options {
		option(cfg)
	}

	w := &waiter{
		waitFuncs:    []WaitFunc{},
		cleanupFuncs: []CleanupFunc{},
	}
	w.ctx, w.cancel = context.WithCancel(cfg.parentCtx)

	if cfg.catchSignals {
		w.ctx, w.cancel = signal.NotifyContext(w.ctx, os.Interrupt, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	}

	return w
}

func (w *waiter) AddWaitFunc(fns ...WaitFunc) {
	w.waitFuncs = append(w.waitFuncs, fns...)
}

func (w *waiter) AddCleanupFunc(fns ...CleanupFunc) {
	w.cleanupFuncs = append(w.cleanupFuncs, fns...)
}

func (w *waiter) Wait() error {
	g, ctx := errgroup.WithContext(w.ctx)

	g.Go(func() error {
		<-ctx.Done()
		w.cancel()
		return nil
	})

	for _, fn := range w.waitFuncs {
		waitFunc := fn

		g.Go(func() error {
			return waitFunc(ctx)
		})
	}

	for _, fn := range w.cleanupFuncs {
		cleanupFunc := fn
		defer cleanupFunc()
	}

	return g.Wait()
}

func (w *waiter) Context() context.Context {
	return w.ctx
}

func (w *waiter) CancelFunc() context.CancelFunc {
	return w.cancel
}
