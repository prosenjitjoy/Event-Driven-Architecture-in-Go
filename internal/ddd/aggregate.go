package ddd

const (
	AggregateIDKey      = "aggregate-id"
	AggregateNameKey    = "aggregate-name"
	AggregateVersionKey = "aggregate-version"
)

type AggregateNamer interface {
	AggregateName() string
}

type Eventer interface {
	AddEvent(string, EventPayload, ...EventOption)
	Events() []AggregateEvent
	ClearEvents()
}

type Aggregate struct {
	Entity
	events []AggregateEvent
}

type AggregateEvent interface {
	Event
	AggregateID() string
	AggregateName() string
	AggregateVersion() int
}

type aggregateEvent struct {
	event
}

var _ interface {
	AggregateNamer
	Eventer
} = (*Aggregate)(nil)

func NewAggregate(id, name string) Aggregate {
	return Aggregate{
		Entity: NewEntity(id, name),
		events: make([]AggregateEvent, 0),
	}
}

func (a Aggregate) AggregateName() string    { return a.name }
func (a Aggregate) Events() []AggregateEvent { return a.events }
func (a *Aggregate) ClearEvents()            { a.events = []AggregateEvent{} }

func (a *Aggregate) AddEvent(name string, payload EventPayload, options ...EventOption) {
	options = append(options, Metadata{
		AggregateIDKey:   a.id,
		AggregateNameKey: a.name,
	})

	a.events = append(a.events, aggregateEvent{
		event: newEvent(name, payload, options...),
	})
}

func (a *Aggregate) setEvents(events []AggregateEvent) { a.events = events }

func (e aggregateEvent) AggregateID() string { return e.metadata.Get(AggregateIDKey).(string) }

func (e aggregateEvent) AggregateName() string { return e.metadata.Get(AggregateNameKey).(string) }

func (e aggregateEvent) AggregateVersion() int { return e.metadata.Get(AggregateVersionKey).(int) }
