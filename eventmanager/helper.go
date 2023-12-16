package eventmanager

import "github.com/google/uuid"

// function to generate a random uuid
func (em *EventManager) GenerateUUID() string {

	// Generate a new random UUID
	uuid := uuid.New()

	return uuid.String()
}
