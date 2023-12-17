// simple event queue. register your subscribers and publish events
package eventmanager

import (
	"fmt"
	"sync"
	"time"

	log "github.com/sirupsen/logrus"
)

func init() {
	log.SetFormatter(&log.JSONFormatter{})
}

// initializes a new event manager
func Initialize(synchronous bool) *EventManager {
	log.Info("initialize a new event manager:")
	log.Info(fmt.Sprintf("- synchronous-mode: %v", synchronous))
	em := &EventManager{
		eventChannel:  make(chan Event),
		synchronous:   synchronous,
		eventHandlers: make([]EventHandler, 0),
	}
	return em
}

// function to publish a new event on eventmanager
func (em *EventManager) Publish(event Event) {

	// set created timestamp
	event.Metadata.UID = em.GenerateUUID()
	event.Metadata.CreatedAt = time.Now()
	log.Info(fmt.Sprintf("publish new event: %s", event.Metadata.UID))

	if em.synchronous {
		em.handleEventSynchronously(event)
	} else {
		go em.handleEventAsynchronously(event)
	}
}

// process an event synchronously
// - send to all handlers in series
// - wait for each handler to finish before sending next event
func (em *EventManager) handleEventSynchronously(event Event) {
	for i := range em.eventHandlers {
		log.Info(fmt.Sprintf("send event %s to handler %s", event.Metadata.UID, em.eventHandlers[i].uid))
		em.eventHandlers[i].handler(event)
		log.Info(fmt.Sprintf("handler %s finished processing event %s", em.eventHandlers[i].uid, event.Metadata.UID))
	}

	// set the ended timestamp
	event.Metadata.EndedAt = time.Now()

	log.Info(fmt.Sprintf("event %s has been processed by all handlers", event.Metadata.UID))

	// push event to history
	em.addEventToHistory(event)
}

// process an event asynchronously
// - send to all handlers in parallel
// - wait for all to finish
func (em *EventManager) handleEventAsynchronously(event Event) {
	var wg sync.WaitGroup

	for i := range em.eventHandlers {
		wg.Add(1)
		go func(index int) {
			log.Info(fmt.Sprintf("send event %s to handler %s", event.Metadata.UID, em.eventHandlers[index].uid))
			defer wg.Done()
			em.eventHandlers[index].handler(event)
			log.Info(fmt.Sprintf("handler %s finished processing event %s", em.eventHandlers[index].uid, event.Metadata.UID))
		}(i)
	}

	// wait for completion of all handlers
	wg.Wait()

	// set the ended timestamp
	event.Metadata.EndedAt = time.Now()

	log.Info(fmt.Sprintf("event %s has been processed by all handlers", event.Metadata.UID))

	// push event to history
	em.addEventToHistory(event)
}

// function to subscribe to events from eventmanager
func (em *EventManager) Subscribe(handler EventHandler) {
	em.eventHandlers = append(em.eventHandlers, handler)
}

// create a new event object
func (em *EventManager) Event(payload []byte) Event {
	return Event{
		Payload: payload,
	}
}
