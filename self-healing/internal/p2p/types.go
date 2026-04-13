package p2p

import (
	"time"
)

// HealthStatus represents the health information shared between nodes
type HealthStatus struct {
	PeerID    string    `json:"peer_id"`
	Hostname  string    `json:"hostname"`
	Status    string    `json:"status"`
	Timestamp time.Time `json:"timestamp"`
}
