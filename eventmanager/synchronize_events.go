package eventmanager

import (
	"fmt"
	"net"

	"google.golang.org/protobuf/proto"
)

// sendEvent sends an Event to a specified address.
func sendEvent(event *Event, address string) error {
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%s", address, "8087"))
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
func receiveEvent(address string) (*Event, error) {
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

	data := make([]byte, 1048576)
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
