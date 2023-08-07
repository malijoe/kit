package events

import (
	"context"
	"sync"
)

type EventHandler interface {
	Handle(ctx context.Context, event Event) error
}

type EventHandlerFunc func(ctx context.Context, event Event) error

func (f EventHandlerFunc) Handle(ctx context.Context, event Event) error {
	return f(ctx, event)
}

type EventSubscriber interface {
	Subscribe(eventType string, handler EventHandler)
}

type EventPublisher interface {
	Publish(ctx context.Context, events ...Event) error
}

type EventDispatcher struct {
	handlers map[string][]EventHandler
	mu       sync.Mutex
}

func NewEventDispatcher() *EventDispatcher {
	return &EventDispatcher{
		handlers: make(map[string][]EventHandler),
	}
}

func (d *EventDispatcher) Subscribe(eventType string, handler EventHandler) {
	d.mu.Lock()
	defer d.mu.Unlock()

	d.handlers[eventType] = append(d.handlers[eventType], handler)
}

func (d *EventDispatcher) Publish(ctx context.Context, events ...Event) (err error) {
	for _, event := range events {
		for _, handler := range d.handlers[event.Type()] {
			err = handler.Handle(ctx, event)
			if err != nil {
				return
			}
		}
	}
	return
}
