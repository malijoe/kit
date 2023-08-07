package events

type Aggregate interface {
	ID() string
	Type() string
	Events() []Event
	Add(event Event) error
	Load([]Event) error
	HasUncommittedEvents() bool
	CommitEvents()
}
