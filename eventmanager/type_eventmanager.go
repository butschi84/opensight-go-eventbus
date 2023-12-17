package eventmanager

// main struct: EventManager
type EventManager struct {
	// channel to process events
	eventChannel chan Event

	// new events will be sent to all registered handlers
	eventHandlers []EventHandler

	// history of last events
	eventHistory []Event

	// should this eventmanager handle events in sync mode (batch processing) or async (parallel processing)
	synchronous bool
}
