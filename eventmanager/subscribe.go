package eventmanager

// function to subscribe to events from eventmanager
func (em *EventManager) Subscribe(handler EventHandler) {
	em.eventHandlers = append(em.eventHandlers, handler)
}
