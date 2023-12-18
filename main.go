package main

import (
	"time"

	"github.com/butschi84/opensight-go-eventbus/eventmanager"
)

func main() {
	// initialize eventmanager
	config := eventmanager.EventManagerConfig{
		MemberListAddress:     "localhost",
		SynchronousProcessing: false, // process evets in parallel
		EventSyncEnabled:      true,  // enable event synchronisation between eventmanager instances
		EventHistoryEnabled:   true,  // enable event history
		EventHistoryLength:    100,   // keep 100 events in history
	}
	em, _ := eventmanager.Initialize(&config)

	time.Sleep(2 * time.Second)

	// add some handlers for testing
	em.Subscribe(em.Handler(handleEvent))
	em.Subscribe(em.Handler(handleEvent2))

	// produce some events for testing
	for range make([]int, 5) {
		em.Publish(em.Event([]byte(`{"test": { "hello": "schmutje" } }`)))
	}

	time.Sleep(3 * time.Second)

	// print history
	em.PrintHistory()
}

func handleEvent(e eventmanager.Event) {
	time.Sleep(1000 * time.Millisecond)
}
func handleEvent2(e eventmanager.Event) {
	time.Sleep(2000 * time.Millisecond)
}
