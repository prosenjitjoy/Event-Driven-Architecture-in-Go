package es

import (
	"context"
	"mall/internal/ddd"
)

type EventSourcedAggregate interface {
	ddd.IDer
	AggregateName() string
	ddd.Eventer
	Versioner
	EventApplier
	EventCommitter
}

type AggregateStore interface {
	Load(ctx context.Context, aggregate EventSourcedAggregate) error
	Save(ctx context.Context, aggregate EventSourcedAggregate) error
}

type AggregateStoreMiddleware func(store AggregateStore) AggregateStore

func AggregateStoreWithMiddleware(store AggregateStore, mws ...AggregateStoreMiddleware) AggregateStore {
	s := store

	// middleware are applied in reverse
	for i := len(mws) - 1; i >= 0; i-- {
		s = mws[i](s)
	}

	return s
}
