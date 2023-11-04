package events

type Aggregate interface {
	ID() string
	Type() string
	Version() int64
	Events() []Event
	Add(event Event) error
	Load([]Event) error
	HasUncommittedEvents() bool
	CommitEvents()
}
