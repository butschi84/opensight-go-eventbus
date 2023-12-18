package eventmanager

import (
	"fmt"
	"sync"

	timestamppb "google.golang.org/protobuf/types/known/timestamppb"
)

// function to publish a new event on eventmanager
func (em *EventManager) Publish(event *Event) {

	// set created timestamp
	event.Metadata.Uid = em.GenerateUUID()
	event.Metadata.CreatedAt = timestamppb.Now()
	logger.Printf("publish new event: %s", event.Metadata.Uid)

	if em.config.SynchronousProcessing {
		em.handleEventSynchronously(*event)
	} else {
		go em.handleEventAsynchronously(*event)
	}

	if em.config.EventSyncEnabled {
		// send event to all members of memberlist
		for _, member := range em.memberList.Members() {
			for i := 0; i < em.config.EventSyncMaxRetransmissions; i++ {
				logger.Printf("send event to member: %s:%d", member.Addr, em.config.EventSyncPort)
				err := em.sendEvent(event, member.Addr.String())
				if err != nil {
					logger.Printf("failed to send event %s to peer %s: %s", event.Metadata.Uid, member.Address(), err.Error())
				} else {
					break
				}
			}
		}
	}
}

// process an event synchronously
// - send to all handlers in series
// - wait for each handler to finish before sending next event
func (em *EventManager) handleEventSynchronously(event Event) {
	for i := range em.eventHandlers {
		logger.Printf("send event %s to handler %s", event.Metadata.Uid, em.eventHandlers[i].uid)
		em.eventHandlers[i].handler(event)
		logger.Printf("handler %s finished processing event %s", em.eventHandlers[i].uid, event.Metadata.Uid)
	}

	// set the ended timestamp
	event.Metadata.EndedAt = timestamppb.Now()

	logger.Printf("event %s has been processed by all handlers", event.Metadata.Uid)

	// push event to history
	if em.config.EventHistoryEnabled {
		em.addEventToHistory(event)
	}
}

// process an event asynchronously
// - send to all handlers in parallel
// - wait for all to finish
func (em *EventManager) handleEventAsynchronously(event Event) {
	var wg sync.WaitGroup

	for i := range em.eventHandlers {
		wg.Add(1)
		go func(index int) {
			logger.Printf("send event %s to handler %s", event.Metadata.Uid, em.eventHandlers[index].uid)
			defer wg.Done()
			em.eventHandlers[index].handler(event)
			logger.Printf("handler %s finished processing event %s", em.eventHandlers[index].uid, event.Metadata.Uid)
		}(i)
	}

	// wait for completion of all handlers
	wg.Wait()

	// set the ended timestamp
	event.Metadata.EndedAt = timestamppb.Now()

	logger.Printf(fmt.Sprintf("event %s has been processed by all handlers", event.Metadata.Uid))

	// push event to history
	if em.config.EventHistoryEnabled {
		em.addEventToHistory(event)
	}
}
