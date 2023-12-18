package eventmanager

import (
	"fmt"
	"time"
)

func (em *EventManager) addEventToHistory(event *Event) {
	// Shift existing events to the right
	n := len(em.eventHistory)
	if n < em.config.EventHistoryLength {
		// If the array is not full, increase its size
		em.eventHistory = append(em.eventHistory, Event{})
		copy(em.eventHistory[1:], em.eventHistory[:n])
	} else {
		// If the array is already full, drop the last event
		copy(em.eventHistory[1:], em.eventHistory[:n-1])
	}
	// Add the new event at the beginning
	em.eventHistory[0] = *event
}

// get last events that were processed
func (em *EventManager) History() []Event {
	return em.eventHistory
}

func (em *EventManager) PrintHistory() {
	fmt.Printf("+----------+--------------------------------------+---------------------+---------------------+----------+--------------+\n")
	fmt.Printf("| id       | event uid                            | created             | ended               | duration | synchronized |\n")
	fmt.Printf("+----------+--------------------------------------+---------------------+---------------------+----------+--------------+\n")
	for i, ev := range em.History() {
		time1 := time.Unix(ev.Metadata.CreatedAt.GetSeconds(), int64(ev.Metadata.CreatedAt.GetNanos()))
		time2 := time.Unix(ev.Metadata.EndedAt.GetSeconds(), int64(ev.Metadata.EndedAt.GetNanos()))
		duration := time2.Sub(time1)
		secondsDifference := int64(duration.Seconds())

		fmt.Printf("| %08d | %s | %s | %s | %08d | %-012v |\n",
			len(em.History())-i,
			ev.Metadata.Uid,
			ev.Metadata.CreatedAt.AsTime().Local().Format("2006-01-02 15:04:05"),
			ev.Metadata.EndedAt.AsTime().Local().Format("2006-01-02 15:04:05"),
			secondsDifference,
			ev.Synchronized,
		)
	}
	fmt.Printf("+----------+--------------------------------------+---------------------+---------------------+----------+--------------+\n")
}
