package eventmanager

func (e *EventManager) addEventToHistory(event Event) {
	// Shift existing events to the right
	n := len(e.eventHistory)
	if n < 100 {
		// If the array is not full, increase its size
		e.eventHistory = append(e.eventHistory, Event{})
		copy(e.eventHistory[1:], e.eventHistory[:n])
	} else {
		// If the array is already full, drop the last event
		copy(e.eventHistory[1:], e.eventHistory[:n-1])
	}
	// Add the new event at the beginning
	e.eventHistory[0] = event
}

// get last events that were processed
func (em *EventManager) History() []Event {
	return em.eventHistory
}
