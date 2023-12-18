package eventmanager

import (
	"fmt"
	"net"

	"google.golang.org/protobuf/proto"
)

// sendEvent sends an Event to a specified address.
func (em *EventManager) sendEvent(event *Event, address string) error {
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", address, em.config.EventSyncPort))
	if err != nil {
		return err
	}
	defer conn.Close()

	data, err := proto.Marshal(event)
	if err != nil {
		return err
	}

	_, err = conn.Write(data)
	return err
}

// receiveEvent listens for incoming events on a specified address.
func (em *EventManager) receiveEvents(address string) {
	logger.Printf("start listening for events on: %s", address)

	ln, err := net.Listen("tcp", address)
	if err != nil {
		logger.Printf("error listening: %v", err)
		return
	}
	defer ln.Close()

	for {
		conn, err := ln.Accept()
		if err != nil {
			logger.Printf("error accepting connection: %v", err)
			ln.Close()
			continue // Retry accepting
		}

		// Handle the connection asynchronously
		go em.handleConnection(conn)
	}
}

func (em *EventManager) handleConnection(conn net.Conn) {
	defer conn.Close()

	data := make([]byte, em.config.EventSyncReceiveBufferSizeBytes)
	n, err := conn.Read(data)
	if err != nil {
		logger.Printf("error reading data: %v", err)
		return
	}

	logger.Printf("received %d bytes", n)

	receivedEvent := &Event{}
	err = proto.Unmarshal(data[:n], receivedEvent)
	if err != nil {
		logger.Printf("error unmarshaling data: %v", err)
		return
	}

	// Process the received event as needed
	em.processEvent(receivedEvent)
}

func (em *EventManager) processEvent(event *Event) {
	// Add your event processing logic here
	logger.Printf("Processing event: %+v", event)
	event.Synchronized = true
	em.Publish(event)
}
