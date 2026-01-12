package domain

import (
	"errors"
	"mall/internal/ddd"
)

var (
	ErrStoreNameIsBlank               = errors.New("the store name cannot be blank")
	ErrStoreLocationIsBlank           = errors.New("the store location cannot be blank")
	ErrStoreIsAlreadyParticipating    = errors.New("the store is already participating")
	ErrStoreIsAlreadyNotParticipating = errors.New("the store is already not participating")
)

type Store struct {
	ddd.AggregateBase
	Name          string
	Location      string
	Participating bool
}

func CreateStore(id, name, location string) (*Store, error) {
	if name == "" {
		return nil, ErrStoreNameIsBlank
	}

	if location == "" {
		return nil, ErrStoreLocationIsBlank
	}

	store := &Store{
		AggregateBase: ddd.AggregateBase{ID: id},
		Name:          name,
		Location:      location,
	}

	store.AddEvent(&StoreCreated{
		Store: store,
	})

	return store, nil
}

func (s *Store) EnableParticipation() error {
	if s.Participating {
		return ErrStoreIsAlreadyParticipating
	}

	s.Participating = true

	s.AddEvent(&StoreParticipationEnabled{
		Store: s,
	})

	return nil
}

func (s *Store) DisableParticipation() error {
	if !s.Participating {
		return ErrStoreIsAlreadyNotParticipating
	}

	s.Participating = false

	s.AddEvent(&StoreParticipationDisabled{
		Store: s,
	})

	return nil
}
