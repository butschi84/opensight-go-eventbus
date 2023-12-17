package main

import (
	"fmt"
	"time"

	"github.com/butschi84/opensight-go-eventbus/eventmanager"
)

func main() {
	em := eventmanager.Initialize(false)

	em.Subscribe(em.Handler(handleEvent))
	em.Subscribe(em.Handler(handleEvent2))

	for range make([]int, 2) {
		em.Publish(em.Event([]byte(`{"test": { "hello": "schmutje" } }`)))
	}

	time.Sleep(3 * time.Second)

	// print history
	history := em.History()
	for i, ev := range history {
		fmt.Printf("event %d: %s\n", len(history)-i, ev.Metadata.UID)
		fmt.Println("-------------------------------------")
		fmt.Println(" - created: " + ev.Metadata.CreatedAt.Local().String())
		fmt.Println(" - ended:   " + ev.Metadata.EndedAt.Local().String())
		fmt.Printf(" - duration:   %f seconds\n", ev.Metadata.EndedAt.Sub(ev.Metadata.CreatedAt).Seconds())
		fmt.Println("")
	}
}

func handleEvent(e eventmanager.Event) {
	time.Sleep(1000 * time.Millisecond)
}
func handleEvent2(e eventmanager.Event) {
	time.Sleep(2000 * time.Millisecond)
}
