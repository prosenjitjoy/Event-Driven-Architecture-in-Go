package domain

import (
	"errors"
	"fmt"
	"mall/internal/ddd"
	"mall/internal/es"
)

const StoreAggregate = "stores.Store"

var (
	ErrStoreNameIsBlank               = errors.New("the store name cannot be blank")
	ErrStoreLocationIsBlank           = errors.New("the store location cannot be blank")
	ErrStoreIsAlreadyParticipating    = errors.New("the store is already participating")
	ErrStoreIsAlreadyNotParticipating = errors.New("the store is already not participating")
)

type Store struct {
	es.Aggregate
	Name          string
	Location      string
	Participating bool
}

var _ interface {
	es.EventApplier
	es.Snapshotter
} = (*Store)(nil)

func NewStore(id string) *Store {
	return &Store{
		Aggregate: es.NewAggregate(id, StoreAggregate),
	}
}

func CreateStore(id, name, location string) (*Store, error) {
	if name == "" {
		return nil, ErrStoreNameIsBlank
	}

	if location == "" {
		return nil, ErrStoreLocationIsBlank
	}

	store := NewStore(id)

	store.AddEvent(StoreCreatedEvent, &StoreCreated{
		Name:     name,
		Location: location,
	})

	return store, nil
}

func (s *Store) EnableParticipation() error {
	if s.Participating {
		return ErrStoreIsAlreadyParticipating
	}

	s.AddEvent(StoreParticipationEnabledEvent, &StoreParticipationToggled{
		Participation: true,
	})

	return nil
}

func (s *Store) DisableParticipation() error {
	if !s.Participating {
		return ErrStoreIsAlreadyNotParticipating
	}

	s.AddEvent(StoreParticipationDisabledEvent, &StoreParticipationToggled{
		Participation: false,
	})

	return nil
}

func (s *Store) Rebrand(name string) error {
	s.AddEvent(StoreRebrandedEvent, &StoreRebranded{
		Name: name,
	})

	return nil
}

func (s *Store) ApplyEvent(event ddd.Event) error {
	switch payload := event.Payload().(type) {
	case *StoreCreated:
		s.Name = payload.Name
		s.Location = payload.Location
	case *StoreParticipationToggled:
		s.Participating = payload.Participation
	case *StoreRebranded:
		s.Name = payload.Name
	default:
		return fmt.Errorf("%T received the event %s with unexpected payload %T", s, event.EventName(), payload)
	}

	return nil
}

func (s *Store) ApplySnapshot(snapshot es.Snapshot) error {
	switch ss := snapshot.(type) {
	case *StoreV1:
		s.Name = ss.Name
		s.Location = ss.Location
		s.Participating = ss.Participation
	default:
		return fmt.Errorf("%T received the unexpected snapshot %T", s, snapshot)
	}

	return nil
}

func (s Store) ToSnapshot() es.Snapshot {
	return StoreV1{
		Name:          s.Name,
		Location:      s.Location,
		Participation: s.Participating,
	}
}
