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

	// load default config
	if config.MemberListBindPort == 0 {
		config.MemberListBindPort = 8081
	}
	if config.EventSyncPort == 0 {
		config.EventSyncPort = 8082
	}

	logger.Printf("initialize a new event manager: %s", nodeName)
	logger.Printf("- memberlist bind port: %d", config.MemberListBindPort)
	logger.Printf("- synchronous-mode:     %v", config.SynchronousProcessing)

	// initialize event manager
	em := &EventManager{
		eventChannel:  make(chan Event),
		config:        config,
		eventHandlers: make([]EventHandler, 0),
		memberList:    nil,
	}

	// initialize memberlist
	err := em.initializeMemberList()
	if err != nil {
		return nil, fmt.Errorf("error while joining memberlist: %s", err.Error())
	}

	// start listening for events, give some time to start listening
	go em.initializeSyncListener()

	return em, nil
}

func (em *EventManager) initializeSyncListener() {
	for {
		receivedEvent, err := receiveEvent(fmt.Sprintf("localhost:%d", em.config.EventSyncPort))
		if err != nil {
			logger.Printf("could not start event receiver / synchronisation: %s", err)
		}
		logger.Printf("Received Event: %+v\n", receivedEvent)
	}
}

// initialize memberlist
func (em *EventManager) initializeMemberList() error {
	logger.Printf("prepare to join memberlist")
	config := memberlist.DefaultLocalConfig()
	config.Name = em.config.name
	config.BindAddr = "127.0.0.1"
	config.BindPort = em.config.MemberListBindPort
	config.Logger = logger
	memberList, err := memberlist.Create(config)

	if err != nil {
		return fmt.Errorf("error initializing cluster node with error %v", err)
	}

	// join the memberlist
	logger.Printf("joining memberlist")
	ma, _ := em.resolveMemberlistDNSName()
	memberList.Join(ma)

	em.memberList = memberList
	return nil
}

func (em *EventManager) Event(event []byte) Event {
	return Event{
		Metadata: &EventMetadata{},
		Payload:  event,
	}
}
