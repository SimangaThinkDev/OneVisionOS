package monitor

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

type TelemetryPayload struct {
	Hostname  string                 `json:"hostname"`
	PeerID    string                 `json:"peer_id"`
	Timestamp time.Time              `json:"timestamp"`
	Metrics   SystemMetrics          `json:"metrics"`
	Status    string                 `json:"status"`
	NIDS      map[string]interface{} `json:"nids"`
}

// SendTelemetry sends system health and security data to the management server
func SendTelemetry(serverURL string, peerID string, metrics SystemMetrics, status string, nidsCount int) {
	hostname, _ := os.Hostname()
	
	payload := TelemetryPayload{
		Hostname:  hostname,
		PeerID:    peerID,
		Timestamp: time.Now(),
		Metrics:   metrics,
		Status:    status,
		NIDS: map[string]interface{}{
			"active":     true,
			"signatures": nidsCount,
		},
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		log.Printf("[TELEMETRY] Failed to marshal payload: %v", err)
		return
	}

	resp, err := http.Post(fmt.Sprintf("%s/api/remote/telemetry", serverURL), "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		log.Printf("[TELEMETRY] Failed to send to %s: %v", serverURL, err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Printf("[TELEMETRY] Server returned error: %s", resp.Status)
	}
}
