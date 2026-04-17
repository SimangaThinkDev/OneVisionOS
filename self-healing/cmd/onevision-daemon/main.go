package main

import (
	"context"
	"encoding/json"
	"encoding/base64"
	"flag"
	"log"
	"github.com/libp2p/go-libp2p/core/peer"
	"net/http"
	"onevision/self-healing/internal/backup"
	"onevision/self-healing/internal/monitor"
	"onevision/self-healing/internal/update"
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
	p2pNode.RegisterCommandProtocol()
	p2pNode.RegisterFileDistProtocol()

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

	// Background loop for P2P health broadcast + telemetry
	go func() {
		ticker := time.NewTicker(30 * time.Second)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				status := wd.GetStatus()
				p2pNode.BroadcastHealthStatus(status)
				monitor.SendTelemetry(
					os.Getenv("MGMT_SERVER_URL"),
					p2pNode.Host.ID().String(),
					monitor.GetSystemMetrics(),
					status,
					len(nidsEngine.Signatures),
				)
			case <-ctx.Done():
				return
			}
		}
	}()

	// Background loop for automatic backups (Phase 19)
	go func() {
		ticker := time.NewTicker(4 * time.Hour)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				log.Println("[BACKUP] Starting scheduled backup...")
				home := os.Getenv("HOME")
				if home == "" {
					home = "/home/student" // Fallback for headless environments
				}
				_, err := backup.CreateSnapshot(backup.BackupConfig{
					SourceDir: home,
					TargetDir: "/var/backups/onevision",
					Excluded:  []string{".cache", "Downloads"},
				})
				if err != nil {
					log.Printf("[BACKUP] Scheduled backup failed: %v", err)
				} else {
					log.Println("[BACKUP] Scheduled backup completed successfully.")
				}
			case <-ctx.Done():
				return
			}
		}
	}()

	// Background loop for automatic updates (Phase 20)
	go func() {
		ticker := time.NewTicker(24 * time.Hour)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				if err := update.RunUpdate(); err != nil {
					log.Printf("[UPDATE] Scheduled update failed: %v", err)
				}
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
		sysMetrics := monitor.GetSystemMetrics()
		secScore := monitor.CalculateSecurityScore(len(nidsEngine.Signatures), 0)

		status := map[string]interface{}{
			"daemon":   "active",
			"watchdog": wd.GetStatus(),
			"metrics":  sysMetrics,
			"p2p": map[string]interface{}{
				"peer_id": p2pNode.Host.ID().String(),
				"peers":   len(p2pNode.Host.Network().Peers()),
				"addrs":   p2pNode.Host.Addrs(),
			},
			"nids": map[string]interface{}{
				"active":         true,
				"signatures":     len(nidsEngine.Signatures),
				"security_score": secScore,
			},
			"repairs": wd.Repairs,
			"time":    time.Now().Format(time.RFC3339),
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(status)
	})

	// Bridge API: Send command to a remote peer via P2P
	http.HandleFunc("/bridge/command", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var req struct {
			PeerID  string   `json:"peer_id"`
			Command string   `json:"command"`
			Args    []string `json:"args"`
		}

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		targetPeer, err := peer.Decode(req.PeerID)
		if err != nil {
			http.Error(w, "Invalid peer_id", http.StatusBadRequest)
			return
		}

		resp, err := p2pNode.SendCommand(targetPeer, p2p.CommandRequest{
			Command: req.Command,
			Args:    req.Args,
		})

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	})

	// Bridge API: Distribute file to a remote peer via P2P
	http.HandleFunc("/bridge/distribute", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var req struct {
			PeerID     string `json:"peer_id"`
			DestPath   string `json:"dest_path"`
			ContentB64 string `json:"content_b64"`
			Mode       uint32 `json:"mode"`
		}

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		content, err := base64.StdEncoding.DecodeString(req.ContentB64)
		if err != nil {
			http.Error(w, "Invalid base64 content", http.StatusBadRequest)
			return
		}

		targetPeer, err := peer.Decode(req.PeerID)
		if err != nil {
			http.Error(w, "Invalid peer_id", http.StatusBadRequest)
			return
		}

		err = p2pNode.DistributeFile(targetPeer, req.DestPath, content, req.Mode)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"status": "success"})
	})

	// Bridge API: Trigger system update (Phase 20)
	http.HandleFunc("/bridge/update", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		go func() {
			if err := update.RunUpdate(); err != nil {
				log.Printf("[UPDATE] Forced update failed: %v", err)
			}
		}()

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"status": "update_started"})
	})

	log.Printf("[API] Internal health check listening on :%s", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatalf("API server failed: %s", err)
	}
}
