package events

type Aggregate interface {
	ID() string
	Type() string
	Version() uint64
	Events() []Event
	Add(event Event) error
	Load([]Event) error
	HasUncommittedEvents() bool
	CommitEvents()
}
