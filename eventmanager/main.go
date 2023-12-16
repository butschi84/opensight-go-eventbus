package eventmanager

import (
	"time"
)

// create a new event manager
func Initialize() *EventManager {
	em := &EventManager{
		eventChannel: make(chan Event),
	}
	return em
}

// function to publish a new event on eventmanager
func (em *EventManager) Publish(event Event) {
	// set created timestamp
	event.UID = em.GenerateUUID()
	event.CreatedAt = time.Now()

	// process the event
	for _, handler := range em.eventHandlers {
		handler(event)
	}

	// set the ended timestamp
	event.EndedAt = time.Now()

	// push event to history
	em.addEventToHistory(event)
}

// function to subscribe to events from eventmanager
func (em *EventManager) Subscribe(handler EventHandler) {
	em.eventHandlers = append(em.eventHandlers, handler)
}
