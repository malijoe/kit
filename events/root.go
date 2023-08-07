package events

const (
	aggregateStartVersion     = -1
	aggregateEventsInitialCap = 10
)

type when func(event Event) error

type AggregateRoot struct {
	id            string
	version       uint64
	globalVersion uint64
	typ           string
	events        []Event
	when          when
}

// interface assertion
var _ Aggregate = (*AggregateRoot)(nil)

func NewAggregateRoot(id string, typ string, when when) (root *AggregateRoot) {
	if when == nil {
		return
	}
	root = &AggregateRoot{
		id:     id,
		typ:    typ,
		when:   when,
		events: make([]Event, 0, aggregateEventsInitialCap),
	}
	return
}

func (r *AggregateRoot) ID() string {
	return r.id
}

func (r *AggregateRoot) Type() string {
	return r.typ
}

func (r *AggregateRoot) Events() (events []Event) {
	// make a copy of the slice to prevent outside manipulation
	events = make([]Event, len(r.events))
	copy(events, r.events)
	return
}

func (r *AggregateRoot) Load(events []Event) (err error) {
	for _, event := range events {
		if event.aggregateType != r.typ {
			return ErrInvalidAggregate
		}

		if err = r.apply(event); err != nil {
			return
		}
		r.version = event.Version()
		r.globalVersion = event.Version()
	}
	return
}

func (r *AggregateRoot) Add(event Event) (err error) {
	if event.aggregateType != r.typ {
		return ErrInvalidAggregate
	}

	if err = r.apply(event); err != nil {
		return
	}
	r.version++
	event.SetVersion(r.version)
	return
}

func (r *AggregateRoot) apply(event Event) (err error) {
	if event.aggregateID != r.id {
		return ErrInvalidAggregateID
	}
	if r.version >= event.version {
		return ErrInvalidEventVersion
	}

	event.aggregateType = r.typ

	if err = r.when(event); err != nil {
		return
	}
	r.events = append(r.events, event)
	return
}

func (r *AggregateRoot) HasUncommittedEvents() bool {
	return len(r.events) > 0
}

func (r *AggregateRoot) CommitEvents() {
	r.globalVersion = r.version
	r.events = make([]Event, 0, aggregateEventsInitialCap)
}
