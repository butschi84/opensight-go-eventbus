package main

import (
	"fmt"
	"time"

	"github.com/butschi84/opensight-go-eventbus/eventmanager"
)

func main() {
	// initialize eventmanager
	config := eventmanager.EventManagerConfig{
		SynchronousProcessing: false,
	}
	em, _ := eventmanager.Initialize(&config)

	time.Sleep(2 * time.Second)

	// add some handlers for testing
	em.Subscribe(em.Handler(handleEvent))
	em.Subscribe(em.Handler(handleEvent2))

	// produce some events for testing
	for range make([]int, 2) {
		em.Publish(em.Event([]byte(`{"test": { "hello": "schmutje" } }`)))
	}

	time.Sleep(3 * time.Second)

	// print history
	history := em.History()
	for i, ev := range history {
		time1 := time.Unix(ev.Metadata.CreatedAt.GetSeconds(), int64(ev.Metadata.CreatedAt.GetNanos()))
		time2 := time.Unix(ev.Metadata.EndedAt.GetSeconds(), int64(ev.Metadata.EndedAt.GetNanos()))
		duration := time2.Sub(time1)
		secondsDifference := int64(duration.Seconds())

		fmt.Printf("event %d: %s\n", len(history)-i, ev.Metadata.Uid)
		fmt.Println("-------------------------------------")
		fmt.Println(" - created: " + ev.Metadata.CreatedAt.String())
		fmt.Println(" - ended:   " + ev.Metadata.EndedAt.String())
		fmt.Printf(" - duration:   %d seconds\n", secondsDifference)
		fmt.Println("")
	}
}

func handleEvent(e eventmanager.Event) {
	time.Sleep(1000 * time.Millisecond)
}
func handleEvent2(e eventmanager.Event) {
	time.Sleep(2000 * time.Millisecond)
}
