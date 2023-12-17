package eventmanager

import "time"

type Event struct {
	Metadata EventMetadata

	// payload of the event
	Payload []byte
}

type EventMetadata struct {
	// generate a random id (uid) for each event
	UID string

	// created timestamp
	CreatedAt time.Time

	// ended timestamp
	EndedAt time.Time
}
