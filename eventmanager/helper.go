package eventmanager

import (
	"fmt"
	"net"

	"github.com/google/uuid"
)

// function to generate a random uuid
func (em *EventManager) GenerateUUID() string {

	// Generate a new random UUID
	uuid := uuid.New()

	return uuid.String()
}

// resolve a dns name to an ip address (string)
func resolveMemberlistDNSName(dnsName string) ([]string, error) {
	ipAddresses, err := net.LookupHost(dnsName)
	if err != nil {
		return []string{}, err
	}

	addresses := make([]string, len(ipAddresses))
	for i, ip := range ipAddresses {
		if ip != "::1" {
			addresses[i] = fmt.Sprintf("%s:8080", ip)
		}
	}
	return addresses, nil
}
