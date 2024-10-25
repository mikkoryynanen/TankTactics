package events

type Event interface {}

type EventHandler func(event Event)

type EventEmitter struct {
	listeners map[string][]EventHandler
}

func NewEventEmiter() *EventEmitter {
	return &EventEmitter{
		listeners: make(map[string][]EventHandler),
	}
}

func (e *EventEmitter) On(eventType string, handler EventHandler) {
	e.listeners[eventType] = append(e.listeners[eventType], handler)
}

func (e *EventEmitter) Emit(eventType string, event Event) {
	if handlers, found := e.listeners[eventType]; found {
		for _, handler := range handlers {
			handler(event)
		}
	}
}