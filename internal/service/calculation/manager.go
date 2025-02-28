package calculation

import (
	"context"
	"sync/atomic"

	"github.com/vilasle/gophermart/internal/logger"
)

type EventType = int

const (
	NewOrder EventType = iota + 1
	NewRule
)

type EventHandler func(context.Context, Event)

type Event struct {
	Type EventType
	Data any
}

type EventManager struct {
	state    atomic.Bool
	events   chan Event
	handlers map[EventType][]EventHandler
}

func NewEventManager(ctx context.Context) *EventManager {
	eventCap := 64
	m := &EventManager{
		events:   make(chan Event, eventCap),
		handlers: make(map[EventType][]EventHandler),
		state:    atomic.Bool{},
	}
	go m.start(ctx)
	return m
}

func (em *EventManager) RaiseEvent(name EventType, data any) {
	if em.Stopped() {
		return
	}

	logger.Debug("got event", "name", name, "data", data)

	em.events <- Event{name, data}
}

func (em *EventManager) Stopped() bool {
	v := em.state.Load()
	return v
}

func (em *EventManager) RegisterHandler(event EventType, handler EventHandler) {
	if _, ok := em.handlers[event]; !ok {
		em.handlers[event] = make([]EventHandler, 0)
	}
	em.handlers[event] = append(em.handlers[event], handler)
}

func (em *EventManager) start(ctx context.Context) {
	worked := true
	limit := make(chan struct{}, 1)
	for worked {
		select {
		case event := <-em.events:
			em.runHandler(ctx, event, limit)
		case <-ctx.Done():
			worked = false
		}
	}
	em.Stop()
}

func (em *EventManager) Stop() {
	if em.Stopped() {
		return
	}
	em.state.Store(true)
	close(em.events)
}

func (em *EventManager) runHandler(ctx context.Context, event Event, limit chan struct{}) {
	handlers, ok := em.handlers[event.Type]
	if !ok {
		return
	}

	for _, handler := range handlers {
		limit <- struct{}{}
		go func(fn EventHandler) { fn(ctx, event); <-limit }(handler)
	}
}
