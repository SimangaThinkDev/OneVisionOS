package monitor

import (
	"log"
	"time"
)

type RepairLog struct {
	Path      string    `json:"path"`
	Timestamp time.Time `json:"timestamp"`
	Status    string    `json:"status"`
}

type Watchdog struct {
	Interval time.Duration
	Status   string
	P2PNode  interface{}
	Repairs  []RepairLog
}

func NewWatchdog(interval time.Duration, p2pNode interface{}) *Watchdog {
	return &Watchdog{
		Interval: interval,
		Status:   "initialized",
		P2PNode:  p2pNode,
	}
}

func (w *Watchdog) Start() {
	w.Status = "running"
	InitializeManifest()

	go func() {
		for {
			log.Println("[Monitor] Checking system integrity...")
			
			for _, file := range Manifest {
				if file.ExpectedHash == "" {
					continue
				}
				
				valid, err := VerifyIntegrity(file.Path, file.ExpectedHash)
				if err != nil {
					log.Printf("[Monitor] ERROR: Could not check %s: %v", file.Path, err)
					continue
				}
				
				if !valid {
					log.Printf("[Monitor] ALERT: Integrity violation detected in %s!", file.Path)
					w.AttemptRecovery(file.Path)
				}
			}
			
			time.Sleep(w.Interval)
		}
	}()
}

func (w *Watchdog) GetStatus() string {
	return w.Status
}
