package events

import "context"

type AggregateStore interface {
	LoadAggregate(ctx context.Context, aggregate Aggregate) error
	SaveAggregate(ctx context.Context, aggregate Aggregate) error
}

type AggregateStoreMiddleware func(store AggregateStore) AggregateStore

func AggregateStoreWithMiddleware(store AggregateStore, middleware ...AggregateStoreMiddleware) AggregateStore {
	s := store
	for _, m := range middleware {
		s = m(s)
	}
	return s
}
