package events

type Event struct {
	Type   EventType
	Object interface{}
}

type EventType string

func NewEvent(eventType EventType, object interface{}) *Event {
	return &Event{eventType, object}
}

type EventHandler func(event *Event)

type EventListener struct {
	Handler EventHandler
}

func NewEventListener(handler EventHandler) *EventListener {
	return &EventListener{handler}
}
