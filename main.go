package main

import (
	"fmt"
	"time"

	"github.com/butschi84/opensight-go-eventbus/eventmanager"
)

func main() {
	em := eventmanager.Initialize()

	em.Subscribe(handleEvent)

	for i := range make([]int, 10) {
		em.Publish(eventmanager.Event{
			Type: eventmanager.EventType(fmt.Sprintf("%d", i)),
		})
	}

	// print history
	history := em.History()
	for i, ev := range history {
		fmt.Printf("event %d: %s\n", len(history)-i, ev.UID)
		fmt.Println("-------------------------------------")
		fmt.Println(" - type:    " + ev.Type)
		fmt.Println(" - created: " + ev.CreatedAt.Local().String())
		fmt.Println(" - ended:   " + ev.EndedAt.Local().String())
		fmt.Printf(" - duration:   %f seconds\n", ev.EndedAt.Sub(ev.CreatedAt).Seconds())
		fmt.Println("")
	}
}

func handleEvent(e eventmanager.Event) {
	time.Sleep(200 * time.Millisecond)
}
