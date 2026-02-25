package es

import (
	"fmt"
	"mall/internal/ddd"
)

type EventApplier interface {
	ApplyEvent(ddd.Event) error
}

type EventCommitter interface {
	CommitEvents()
}

func LoadEvent(v any, event ddd.AggregateEvent) error {
	type loader interface {
		EventApplier
		VersionSetter
	}

	agg, ok := v.(loader)
	if !ok {
		return fmt.Errorf("%T does not have the methods implemented to events", v)
	}

	if err := agg.ApplyEvent(event); err != nil {
		return err
	}

	agg.SetVersion(event.AggregateVersion())

	return nil
}
