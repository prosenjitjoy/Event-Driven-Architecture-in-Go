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

type SagaContext[T any] struct {
	ID           string
	Data         T
	Step         int
	Done         bool
	Compensating bool
}

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

func (s *saga[T]) Name() string            { return s.name }
func (s *saga[T]) ReplyTopic() string      { return s.replyTopic }
func (s *saga[T]) getSteps() []SagaStep[T] { return s.steps }

func (s *SagaContext[T]) advance(steps int) {
	var dir = 1

	if s.Compensating {
		dir = -1
	}

	s.Step = s.Step + (dir * steps)
}

func (s *SagaContext[T]) complete()   { s.Done = true }
func (s *SagaContext[T]) compensate() { s.Compensating = true }
