// simple event queue. register your subscribers and publish events
package eventmanager

import (
	"fmt"
	"log"
	"os"

	petname "github.com/dustinkirkland/golang-petname"
	memberlist "github.com/hashicorp/memberlist"
)

var logger *log.Logger

func init() {
	logger = log.New(
		os.Stderr,
		"eventmanager: ",
		log.Ldate|log.Ltime|log.Lmicroseconds|log.LUTC|log.Lshortfile,
	)
}

// initializes a new event manager
func Initialize(config *EventManagerConfig) (*EventManager, error) {
	nodeName := petname.Generate(2, "-")
	logger.Printf("initialize a new event manager: %s", nodeName)

	// initialize memberlist
	memberList, err := initializeMemberList(nodeName, config.MemberListAddress)
	if err != nil {
		return nil, fmt.Errorf("error while joining memberlist: %s", err.Error())
	}

	// initialize event manager
	logger.Printf("- synchronous-mode: %v", config.SynchronousProcessing)
	em := &EventManager{
		eventChannel:  make(chan Event),
		config:        config,
		eventHandlers: make([]EventHandler, 0),
		memberList:    memberList,
	}

	// start listening for events
	go initializeSyncListener()

	return em, nil
}

func initializeSyncListener() {
	for {
		receivedEvent, err := receiveEvent("localhost:8087")
		if err != nil {
			logger.Printf("could not start event receiver / synchronisation: %s", err)
		}
		logger.Printf("Received Event: %+v\n", receivedEvent)
	}
}

// initialize memberlist
func initializeMemberList(nodeName string, memberListAddress string) (*memberlist.Memberlist, error) {
	logger.Printf("prepare to join memberlist")
	config := memberlist.DefaultLocalConfig()
	config.Name = nodeName
	config.BindAddr = "127.0.0.1"
	config.BindPort = 8080
	config.Logger = logger
	memberList, err := memberlist.Create(config)

	if err != nil {
		return nil, fmt.Errorf("error initializing cluster node with error %v", err)
	}

	// join the memberlist
	logger.Printf("joining memberlist")
	ma, _ := resolveMemberlistDNSName(memberListAddress)
	memberList.Join(ma)

	return memberList, nil
}

func (em *EventManager) Event(event []byte) Event {
	return Event{
		Metadata: &EventMetadata{},
		Payload:  event,
	}
}
