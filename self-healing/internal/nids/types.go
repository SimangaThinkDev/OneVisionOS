package nids

import (
	"time"
)

type AlertLevel string

const (
	LevelLow      AlertLevel = "LOW"
	LevelMedium   AlertLevel = "MEDIUM"
	LevelHigh     AlertLevel = "HIGH"
	LevelCritical AlertLevel = "CRITICAL"
)

type Signature struct {
	ID          string     `json:"id"`
	Name        string     `json:"name"`
	Description string     `json:"description"`
	Level       AlertLevel `json:"level"`
	// Basic matching criteria (could be more complex like Snort rules)
	PayloadPattern string `json:"payload_pattern"`
	DestPort       int    `json:"dest_port"`
}

type Alert struct {
	Timestamp time.Time  `json:"timestamp"`
	Signature Signature  `json:"signature"`
	SourceIP  string     `json:"source_ip"`
	DestIP    string     `json:"dest_ip"`
	Payload   string     `json:"payload_snippet"`
}
