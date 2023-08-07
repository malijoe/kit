package events

import "errors"

var (
	ErrInvalidAggregate    = errors.New("invalid aggregate")
	ErrInvalidAggregateID  = errors.New("invalid aggregate id")
	ErrInvalidEventType    = errors.New("invalid event type")
	ErrInvalidEventVersion = errors.New("invalid event version")
)
