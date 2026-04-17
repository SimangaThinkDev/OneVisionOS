package backup

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"time"
)

type BackupConfig struct {
	SourceDir string
	TargetDir string
	Excluded  []string
}

// CreateSnapshot performs an rsync-based snapshot of the source directory
func CreateSnapshot(config BackupConfig) (string, error) {
	if _, err := os.Stat(config.TargetDir); os.IsNotExist(err) {
		if err := os.MkdirAll(config.TargetDir, 0755); err != nil {
			return "", fmt.Errorf("failed to create target dir: %w", err)
		}
	}

	timestamp := time.Now().Format("20060102-150405")
	snapshotPath := filepath.Join(config.TargetDir, "snapshot-"+timestamp)

	args := []string{"-az", "--delete"}
	for _, ex := range config.Excluded {
		args = append(args, "--exclude", ex)
	}
	
	// Link to latest snapshot for incremental efficiency (if it exists)
	latest := filepath.Join(config.TargetDir, "latest")
	if _, err := os.Stat(latest); err == nil {
		args = append(args, "--link-dest", latest)
	}

	args = append(args, config.SourceDir+"/", snapshotPath)

	cmd := exec.Command("rsync", args...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return string(output), fmt.Errorf("rsync failed: %w", err)
	}

	// Update 'latest' symlink
	os.Remove(latest)
	if err := os.Symlink(snapshotPath, latest); err != nil {
		log.Printf("[BACKUP] Warning: failed to update 'latest' symlink: %v", err)
	}

	return snapshotPath, nil
}

// RestoreSnapshot restores from a specific snapshot path
func RestoreSnapshot(snapshotPath string, targetDir string) (string, error) {
	if _, err := os.Stat(snapshotPath); os.IsNotExist(err) {
		return "", fmt.Errorf("snapshot does not exist: %s", snapshotPath)
	}

	args := []string{"-az", "--delete", snapshotPath+"/", targetDir}
	cmd := exec.Command("rsync", args...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return string(output), fmt.Errorf("restore failed: %w", err)
	}

	return string(output), nil
}

// ListSnapshots returns a list of available snapshot directory names
func ListSnapshots(targetDir string) ([]string, error) {
	entries, err := os.ReadDir(targetDir)
	if err != nil {
		return nil, err
	}

	var snapshots []string
	for _, e := range entries {
		if e.IsDir() && e.Name() != "latest" {
			snapshots = append(snapshots, e.Name())
		}
	}
	return snapshots, nil
}
