package events

import (
	"cmp"
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
	version       int64
	metadata      []byte
}

func EventSortFunc(a, b Event) int {
	return cmp.Compare(a.version, b.version)
}

func (event Event) Marshal() any {
	return struct {
		Id            string    `json:"id" yaml:"id"`
		Type          string    `json:"type" yaml:"type"`
		Data          []byte    `json:"data,omitempty" yaml:"data,omitempty"`
		Timestamp     time.Time `json:"timestamp" yaml:"timestamp"`
		AggregateId   string    `json:"aggregateId" yaml:"aggregateId"`
		AggregateType string    `json:"aggregateType" yaml:"aggregateType"`
		Version       int64     `json:"version" yaml:"version"`
		Metadata      []byte    `json:"metadata,omitempty" yaml:"metadata,omitempty"`
	}{
		Id:            event.id,
		Type:          event.typ,
		Data:          event.data,
		Timestamp:     event.timestamp,
		AggregateId:   event.aggregateID,
		AggregateType: event.aggregateType,
		Version:       event.version,
		Metadata:      event.metadata,
	}
}

func (event Event) MarshalJSON() ([]byte, error) {
	return json.Marshal(event.Marshal())
}

func (event *Event) Unmarshal(unmarshal func(any) error) error {
	var obj struct {
		Id            string    `json:"id" yaml:"id"`
		Type          string    `json:"type" yaml:"type"`
		Data          []byte    `json:"data" yaml:"data"`
		Timestamp     time.Time `json:"timestamp" yaml:"timestamp"`
		AggregateId   string    `json:"aggregateId" yaml:"aggregateId"`
		AggregateType string    `json:"aggregateType" yaml:"aggregateType"`
		Version       int64     `json:"version" yaml:"version"`
		Metadata      []byte    `json:"metadata" yaml:"metadata"`
	}
	if err := unmarshal(&obj); err != nil {
		return err
	}
	event.id = obj.Id
	event.typ = obj.Type
	event.data = obj.Data
	event.timestamp = obj.Timestamp
	event.aggregateID = obj.AggregateId
	event.aggregateType = obj.AggregateType
	event.version = obj.Version
	event.metadata = obj.Metadata
	return nil
}

func (event *Event) UnmarshalJSON(data []byte) error {
	return event.Unmarshal(func(obj any) error {
		return json.Unmarshal(data, obj)
	})
}

func NewEvent(root Aggregate, typ string) Event {
	return Event{
		id:            uuid.NewV4().String(),
		aggregateType: root.Type(),
		aggregateID:   root.ID(),
		version:       root.Version() + 1,
		typ:           typ,
		timestamp:     time.Now().UTC(),
	}
}

func (e Event) ID() string {
	return e.id
}

func (e *Event) SetID(id string) {
	e.id = id
}

func (e Event) Type() string {
	return e.typ
}

func (e *Event) SetType(typ string) {
	e.typ = typ
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

func (e Event) Version() int64 {
	return e.version
}

func (e *Event) SetVersion(version int64) {
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
