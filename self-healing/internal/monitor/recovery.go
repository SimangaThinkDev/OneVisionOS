package monitor

import (
	"context"
	"log"
	"time"
)

// RecoveryNode defines the interface for fetching clean files from peers
type RecoveryNode interface {
	FetchFileFromAnyPeer(ctx context.Context, path string, targetPath string) error
}

func (w *Watchdog) AttemptRecovery(filePath string) {
	log.Printf("[Recovery] Attempting to recover %s...", filePath)

	repair := RepairLog{
		Path:      filePath,
		Timestamp: time.Now(),
		Status:    "FAILED",
	}

	if w.P2PNode == nil {
		log.Printf("[Recovery] Error: No P2P node available for recovery")
		w.Repairs = append(w.Repairs, repair)
		return
	}

	rn, ok := w.P2PNode.(RecoveryNode)
	if !ok {
		log.Printf("[Recovery] Error: P2P node does not implement RecoveryNode interface")
		w.Repairs = append(w.Repairs, repair)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	err := rn.FetchFileFromAnyPeer(ctx, filePath, filePath)
	if err != nil {
		log.Printf("[Recovery] Failed to recover %s: %v", filePath, err)
		w.Repairs = append(w.Repairs, repair)
		return
	}

	log.Printf("[Recovery] Successfully restored %s from peer network", filePath)
	repair.Status = "RESTORED"
	w.Repairs = append(w.Repairs, repair)
}
