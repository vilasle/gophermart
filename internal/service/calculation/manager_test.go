package calculation

import (
	"context"
	"reflect"
	"sync"
	"testing"
	"time"
)

func TestEventManager(t *testing.T) {

	baseCtx := context.Background()

	type fields struct {
		events []Event
	}

	type result struct {
		mx     *sync.Mutex
		events []Event
	}

	tests := []struct {
		name string
		fields
		result
	}{
		{
			name: "test event manager",
			fields: fields{
				events: []Event{
					{
						Type: NewOrder,
						Data: "new order",
					},
					{
						Type: NewRule,
						Data: "new rule",
					},
				},
			},
			result: result{
				mx:     &sync.Mutex{},
				events: make([]Event, 0),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fn := func(ctx context.Context, event Event) {
				tt.result.mx.Lock()
				defer tt.result.mx.Unlock()
				tt.result.events = append(tt.result.events, event)
			}
			ctx, cancel := context.WithCancel(baseCtx)
			//check starting manager
			em := NewEventManager()
			em.Start(ctx)

			//check registration
			em.RegisterHandler(NewOrder, fn)
			em.RegisterHandler(NewRule, fn)

			//check raising events
			for _, event := range tt.fields.events {
				em.RaiseEvent(event.Type, event.Data)
			}

			//raise event which does not have subscribers
			em.RaiseEvent(42, "wrong event")

			time.Sleep(time.Millisecond * 500)
			ok := reflect.DeepEqual(tt.fields.events, tt.result.events)
			if !ok {
				t.Errorf("has different has and got events")
			}

			currentLen := len(tt.result.events)
			//check stopping manager
			cancel()

			time.Sleep(time.Millisecond * 500)
			//did not add new event
			for i := 0; i < 10; i++ {
				em.RaiseEvent(NewOrder, "new order")
			}
			time.Sleep(time.Millisecond * 500)

			if currentLen != len(tt.result.events) {
				t.Errorf("added event after stopping manager")
			}
		})
	}
}
