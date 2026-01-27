package sec

import "mall/internal/am"

const (
	SagaCommandIDHeader   = am.CommandHeaderPrefix + "SAGA_ID"
	SagaCommandNameHeader = am.CommandHeaderPrefix + "SAGA_NAME"

	SagaReplyIDHeader   = am.ReplyHeaderPrefix + "SAGA_ID"
	SagaReplyNameHeader = am.ReplyHeaderPrefix + "SAGA_NAME"
)

const (
	isCompensating  = true
	notCompensating = false
)

type Saga[T any] interface {
	AddStep() SagaStep[T]
	Name() string
	ReplyTopic() string
	getSteps() []SagaStep[T]
}

type saga[T any] struct {
	name       string
	replyTopic string
	steps      []SagaStep[T]
}

func NewSaga[T any](name, replyTopic string) Saga[T] {
	return &saga[T]{
		name:       name,
		replyTopic: replyTopic,
	}
}

func (s *saga[T]) AddStep() SagaStep[T] {
	step := &sagaStep[T]{
		actions: map[bool]StepActionFunc[T]{
			isCompensating:  nil,
			notCompensating: nil,
		},
		handlers: map[bool]map[string]StepReplyHandlerFunc[T]{
			isCompensating:  {},
			notCompensating: {},
		},
	}

	s.steps = append(s.steps, step)

	return step
}

func (s *saga[T]) Name() string {
	return s.name
}

func (s *saga[T]) ReplyTopic() string {
	return s.replyTopic
}

func (s *saga[T]) getSteps() []SagaStep[T] {
	return s.steps
}
