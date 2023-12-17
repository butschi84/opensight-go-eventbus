package eventmanager

// eventhandler
type EventHandler struct {
	uid     string
	handler func(event Event)
}

// function to create a new event handler
func (em *EventManager) Handler(handlerFunction func(event Event)) EventHandler {
	return EventHandler{
		uid:     em.GenerateUUID(),
		handler: handlerFunction,
	}
}
