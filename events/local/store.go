package local

import (
	"context"
	"encoding/json"
	"slices"

	"github.com/boltdb/bolt"
	"github.com/malijoe/kit/errors"
	"github.com/malijoe/kit/events"
)

const (
	collection = "EVENTS"
	// terminator is the null character and is used as a terminator
	terminator = "\x00"
)

func toEntryKey(key string) []byte {
	return []byte(key + terminator)
}

type store struct {
	db *bolt.DB
}

func NewLocalStore(db *bolt.DB) *store {
	return &store{db: db}
}

func (s *store) SaveAggregate(ctx context.Context, agg events.Aggregate) error {
	txn, err := s.db.Begin(true)
	if err != nil {
		return err
	}
	defer txn.Rollback()

	typeBckt, err := txn.CreateBucketIfNotExists([]byte(agg.Type()))
	if err != nil {
		return err
	}
	aggBckt, err := typeBckt.CreateBucketIfNotExists([]byte(agg.ID()))
	if err != nil {
		return err
	}

	eventData := make(map[string][]byte)
	stream := agg.Events()
	for _, event := range stream {
		data, err := json.Marshal(&event)
		if err != nil {
			return err
		}
		eventData[event.ID()] = data
	}

	for id, data := range eventData {
		if err := aggBckt.Put(toEntryKey(id), data); err != nil {
			return err
		}
	}

	return txn.Commit()
}

func (s *store) LoadAggregate(ctx context.Context, agg events.Aggregate) error {
	txn, err := s.db.Begin(false)
	if err != nil {
		return err
	}
	typeBckt := txn.Bucket([]byte(agg.Type()))
	if typeBckt == nil {
		return errors.NotFoundErrorf("did not find events for %s aggregate", agg.Type())
	}
	aggBckt := typeBckt.Bucket([]byte(agg.ID()))
	if aggBckt == nil {
		return errors.NotFoundErrorf("did not find events for %s aggregate with id %s", agg.Type(), agg.ID())
	}
	var stream []events.Event
	cursor := aggBckt.Cursor()
	for k, v := cursor.First(); k != nil; k, v = cursor.Next() {
		var event events.Event
		if err := json.Unmarshal(v, &event); err != nil {
			return err
		}
		stream = append(stream, event)
	}

	slices.SortFunc(stream, events.EventSortFunc)
	return agg.Load(stream)
}

func (s *store) GetAggregatesOfType(ctx context.Context, typ string) (streams map[string][]events.Event, _ error) {
	txn, err := s.db.Begin(false)
	if err != nil {
		return nil, err
	}

	typeBckt := txn.Bucket([]byte(typ))
	if typeBckt == nil {
		return nil, errors.NotFoundErrorf("did not find events for %s aggregate", typ)
	}
	streams = make(map[string][]events.Event)
	cursor := typeBckt.Cursor()
	for k, _ := cursor.First(); k != nil; k, _ = cursor.Next() {
		// iterate over the individual aggregates
		// key value should be aggregate id
		aggBckt := typeBckt.Bucket(k)
		if aggBckt == nil {
			// go next if there was no sub bckt
			continue
		}

		cursor := aggBckt.Cursor()
		var aggEvents []events.Event
		for k, v := cursor.First(); k != nil; k, v = cursor.Next() {
			// key can be ignored as long as it's not nil
			// event id is captured in value data
			var event events.Event
			if err := json.Unmarshal(v, &event); err != nil {
				return nil, err
			}
			aggEvents = append(aggEvents, event)
		}
		streams[string(k)] = aggEvents
	}
	return streams, nil
}
