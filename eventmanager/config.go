package eventmanager

type EventManagerConfig struct {
	name string

	MemberListAddress  string
	MemberListBindPort int

	EventSyncEnabled                bool
	EventSyncPort                   int
	EventSyncReceiveBufferSizeBytes int

	SynchronousProcessing bool
}
