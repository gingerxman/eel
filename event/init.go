package event

import (
	"context"
	"time"
)

type Event struct {
	Name string
	Tag  string
}

func NewEvent(name, tag string) *Event {
	event := new(Event)
	event.Name = name
	event.Tag = tag
	return event
}

func Emit(ctx context.Context, e *Event, data map[string]interface{}) {

	eventData := map[string]interface{}{
		"_event_name": e.Name,
		"_time":       time.Now().Format("2006-01-02 15:04:05"),
		"data":        data,
	}

	for _, validEngine := range getValidEngines() {
		validEngine.Send(ctx, eventData, e.Tag)
	}
}
