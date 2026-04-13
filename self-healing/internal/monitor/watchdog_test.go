package monitor

import (
	"testing"
	"time"
)

func TestWatchdogStatus(t *testing.T) {
	w := NewWatchdog(1*time.Second, nil)
	if w.GetStatus() != "initialized" {
		t.Errorf("expected initialized, got %s", w.GetStatus())
	}

	w.Start()
	// Allow a tiny bit of time for the goroutine to update status if it were complex,
	// but here it's synchronous in Start().
	if w.GetStatus() != "running" {
		t.Errorf("expected running, got %s", w.GetStatus())
	}
}
