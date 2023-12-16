package eventmanager

import "time"

type Event struct {
	// generate a random id (uid) for each event
	UID string

	// Data is the payload of the event
	Data interface{}

	// event type
	Type EventType

	// event description (optional)
	Description string

	// created timestamp
	CreatedAt time.Time

	// ended timestamp
	EndedAt time.Time
}
