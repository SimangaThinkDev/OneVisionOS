package monitor

import (
	"log"
	"time"
)

type Watchdog struct {
	Interval time.Duration
	Status   string
}

func NewWatchdog(interval time.Duration) *Watchdog {
	return &Watchdog{
		Interval: interval,
		Status:   "initialized",
	}
}

func (w *Watchdog) Start() {
	w.Status = "running"
	go func() {
		for {
			log.Println("[Monitor] Checking system integrity...")
			// TODO: Implement actual integrity checks here in Phase 12
			time.Sleep(w.Interval)
		}
	}()
}

func (w *Watchdog) GetStatus() string {
	return w.Status
}
