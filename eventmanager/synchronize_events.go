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
func (em *EventManager) receiveEvent(address string) (*Event, error) {
	logger.Printf("start listening for events on: %s", address)

	ln, err := net.Listen("tcp", address)
	if err != nil {
		return nil, err
	}
	defer ln.Close()

	conn, err := ln.Accept()
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	data := make([]byte, em.config.EventSyncReceiveBufferSizeBytes)
	n, err := conn.Read(data)
	if err != nil {
		return nil, err
	}
	logger.Printf("received %d bytes", n)

	receivedEvent := &Event{}
	err = proto.Unmarshal(data[:n], receivedEvent)
	if err != nil {
		return nil, err
	}

	return receivedEvent, nil
}
