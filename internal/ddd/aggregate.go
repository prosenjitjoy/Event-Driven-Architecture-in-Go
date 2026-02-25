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
	GetEvents() []AggregateEvent
	ClearEvents()
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

type Aggregate interface {
	IDer
	AggregateNamer
	Eventer
	IDSetter
	NameSetter
}

type aggregate struct {
	Entity
	events []AggregateEvent
}

var _ Aggregate = (*aggregate)(nil)

func NewAggregate(id, name string) *aggregate {
	return &aggregate{
		Entity: NewEntity(id, name),
		events: make([]AggregateEvent, 0),
	}
}

func (a aggregate) AggregateName() string       { return a.EntityName() }
func (a aggregate) GetEvents() []AggregateEvent { return a.events }
func (a *aggregate) ClearEvents()               { a.events = []AggregateEvent{} }

func (a *aggregate) AddEvent(name string, payload EventPayload, options ...EventOption) {
	options = append(options, Metadata{
		AggregateIDKey:   a.ID(),
		AggregateNameKey: a.EntityName(),
	})

	a.events = append(a.events, aggregateEvent{
		event: newEvent(name, payload, options...),
	})
}

func (a *aggregate) setEvents(events []AggregateEvent) { a.events = events }

func (e aggregateEvent) AggregateID() string {
	aggregateId := e.metadata.Get(AggregateIDKey)

	return aggregateId.(string)
}

func (e aggregateEvent) AggregateName() string {
	aggregateName := e.metadata.Get(AggregateNameKey)

	return aggregateName.(string)
}

func (e aggregateEvent) AggregateVersion() int {
	aggregateVersion := e.metadata.Get(AggregateVersionKey)

	return aggregateVersion.(int)
}
