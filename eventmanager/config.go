package eventmanager

type EventManagerConfig struct {
	name string

	MemberListAddress  string
	MemberListBindPort int

	EventSyncPort int

	SynchronousProcessing bool
}
