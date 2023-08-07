package events

import (
	"encoding/json"
	"time"

	uuid "github.com/satori/go.uuid"
)

type Event struct {
	id            string
	typ           string
	data          []byte
	timestamp     time.Time
	aggregateType string
	aggregateID   string
	version       uint64
	metadata      []byte
}

func NewEvent(root Aggregate, typ string) Event {
	return Event{
		id:            uuid.NewV4().String(),
		aggregateType: root.Type(),
		aggregateID:   root.ID(),
		typ:           typ,
		timestamp:     time.Now().UTC(),
	}
}

func (e Event) ID() string {
	return e.id
}

func (e Event) Type() string {
	return e.typ
}

func (e Event) Data() []byte {
	return e.data
}

func (e Event) Timestamp() time.Time {
	return e.timestamp
}

func (e Event) Metadata() []byte {
	return e.metadata
}

func (e Event) AggregateType() string {
	return e.aggregateType
}

func (e Event) AggregateID() string {
	return e.aggregateID
}

func (e Event) Version() uint64 {
	return e.version
}

func (e Event) SetVersion(version uint64) {
	e.version = version
}

func (e *Event) SetData(data []byte) {
	e.data = data
}

func (e *Event) GetJSONData(data any) error {
	return json.Unmarshal(e.data, data)
}

func (e *Event) SetJSONData(data any) (err error) {
	var dataBytes []byte
	dataBytes, err = json.Marshal(data)
	if err != nil {
		return
	}

	e.data = dataBytes
	return
}
