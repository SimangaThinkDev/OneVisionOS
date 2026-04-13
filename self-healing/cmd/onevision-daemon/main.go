package main

import (
	"context"
	"encoding/json"
	"flag"
	"log"
	"net/http"
	"onevision/self-healing/internal/monitor"
	"onevision/self-healing/internal/nids"
	"onevision/self-healing/internal/p2p"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	apiPort := flag.String("port", "8081", "API port for health check")
	flag.Parse()

	log.Println("Starting OneVisionOS Self-Healing Daemon...")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Initialize P2P Node
	p2pNode, err := p2p.NewNode(ctx)
	if err != nil {
		log.Fatalf("Failed to initialize P2P node: %v", err)
	}

	if err := p2pNode.SetupDiscovery(); err != nil {
		log.Printf("[P2P] Discovery setup failed: %v", err)
	}

	// Register protocols
	p2pNode.RegisterFileSyncProtocol()

	// Initialize NIDS
	nidsEngine := nids.NewEngine()
	nidsEngine.LoadSignatures(nids.DefaultSignatures())
	nidsEngine.StartAutoCapture()

	// Initialize Watchdog
	wd := monitor.NewWatchdog(10*time.Second, p2pNode)
	wd.Start()

	// Setup Health Check API
	go setupHealthCheckAPI(wd, p2pNode, nidsEngine, *apiPort)

	// Listen for NIDS alerts in background
	go func() {
		for alert := range nidsEngine.Alerts {
			log.Printf("[SECURITY] %s: Found %s from %s to %s", alert.Signature.Level, alert.Signature.Name, alert.SourceIP, alert.DestIP)
			// In Phase 14, we will add automated isolation logic here
		}
	}()

	// Background loop for P2P health broadcast
	go func() {
		ticker := time.NewTicker(30 * time.Second)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				status := wd.GetStatus()
				p2pNode.BroadcastHealthStatus(status)
			case <-ctx.Done():
				return
			}
		}
	}()

	// Listen for signals
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	log.Println("Daemon is active. Press Ctrl+C to stop.")

	// Wait for a termination signal
	sig := <-sigs
	log.Printf("Received signal: %s. Shutting down gracefully...", sig)
}

func setupHealthCheckAPI(wd *monitor.Watchdog, p2pNode *p2p.Node, nidsEngine *nids.Engine, port string) {
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		status := map[string]interface{}{
			"daemon":   "active",
			"watchdog": wd.GetStatus(),
			"p2p": map[string]interface{}{
				"peer_id": p2pNode.Host.ID().String(),
				"peers":   len(p2pNode.Host.Network().Peers()),
				"addrs":   p2pNode.Host.Addrs(),
			},
			"nids": map[string]interface{}{
				"active":     true,
				"signatures": len(nidsEngine.Signatures),
			},
			"time": time.Now().Format(time.RFC3339),
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(status)
	})

	log.Printf("[API] Internal health check listening on :%s", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatalf("API server failed: %s", err)
	}
}
