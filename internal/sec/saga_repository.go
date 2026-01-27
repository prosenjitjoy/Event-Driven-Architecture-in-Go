package sec

import (
	"context"
	"fmt"
	"mall/internal/registry"
)

type SagaStore interface {
	Load(ctx context.Context, sagaName, sagaID string) (*SagaContext[[]byte], error)
	Save(ctx context.Context, sagaName string, sagaCtx *SagaContext[[]byte]) error
}

type SagaRepository[T any] struct {
	reg   registry.Registry
	store SagaStore
}

func NewSagaRepository[T any](reg registry.Registry, store SagaStore) SagaRepository[T] {
	return SagaRepository[T]{
		reg:   reg,
		store: store,
	}
}

func (r SagaRepository[T]) Load(ctx context.Context, sagaName, sagaID string) (*SagaContext[T], error) {
	byteCtx, err := r.store.Load(ctx, sagaName, sagaID)
	if err != nil {
		return nil, err
	}

	v, err := r.reg.Deserialize(sagaName, byteCtx.Data)
	if err != nil {
		return nil, err
	}

	data, ok := v.(T)
	if !ok {
		return nil, fmt.Errorf("%T is not the expected type %T", v, data)
	}

	return &SagaContext[T]{
		ID:           byteCtx.ID,
		Data:         data,
		Step:         byteCtx.Step,
		Done:         byteCtx.Done,
		Compensating: byteCtx.Compensating,
	}, nil
}

func (r SagaRepository[T]) Save(ctx context.Context, sagaName string, sagaCtx *SagaContext[T]) error {
	data, err := r.reg.Serialize(sagaName, sagaCtx.Data)
	if err != nil {
		return err
	}

	return r.store.Save(ctx, sagaName, &SagaContext[[]byte]{
		ID:           sagaCtx.ID,
		Data:         data,
		Step:         sagaCtx.Step,
		Done:         sagaCtx.Done,
		Compensating: sagaCtx.Compensating,
	})
}
