package eventmanager

type EventManagerConfig struct {
	name string

	MemberListAddress  string
	MemberListBindPort int

	EventSyncEnabled                bool
	EventSyncPort                   int
	EventSyncReceiveBufferSizeBytes int
	EventSyncMaxRetransmissions     int

	SynchronousProcessing bool

	EventHistoryEnabled bool

	// number of last events that should be kept in history
	EventHistoryLength int
}
