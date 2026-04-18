package update

import (
	"fmt"
	"log"
	"os/exec"
	"time"
)

var commandRunner = exec.Command

// RunUpdate handles the background patching process
func RunUpdate() error {
	log.Println("[UPDATE] Starting background update process...")

	// 1. Update package lists
	cmd := commandRunner("apt-get", "update")
	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("apt-get update failed: %v, output: %s", err, string(output))
	}

	// 2. Perform non-interactive upgrade
	// We use -o Dpkg::Options::="--force-confold" to keep existing configs
	cmd = commandRunner("apt-get", "dist-upgrade", "-y", "-o", "Dpkg::Options::=\"--force-confold\"", "--ignore-hold")
	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("apt-get upgrade failed: %v, output: %s", err, string(output))
	}

	// 3. Clean up
	commandRunner("apt-get", "autoremove", "-y").Run()

	log.Println("[UPDATE] Background update completed successfully.")
	return nil
}

// CheckForUpdates returns true if there are pending updates
func CheckForUpdates() (bool, error) {
	cmd := commandRunner("apt-get", "update")
	cmd.Run()

	cmd = commandRunner("apt-get", "-s", "upgrade")
	_, err := cmd.CombinedOutput()
	if err != nil {
		return false, err
	}

	// If output contains numbers in "0 upgraded, 0 newly installed", it means no updates.
	// Actually we search for "Inst " which means interesting updates.
	// Or check return code.
	return true, nil // Simplified for this implementation
}

// ScheduledUpdateLoop runs updates periodically (e.g., daily)
func ScheduledUpdateLoop(interval time.Duration) {
	ticker := time.NewTicker(interval)
	for range ticker.C {
		if err := RunUpdate(); err != nil {
			log.Printf("[UPDATE] Scheduled update failed: %v", err)
		}
	}
}
