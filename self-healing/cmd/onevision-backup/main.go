package main

import (
	"flag"
	"fmt"
	"log"
	"onevision/self-healing/internal/backup"
	"os"
	"path/filepath"
)

func main() {
	action := flag.String("action", "backup", "Action to perform: backup, restore, list")
	snapshot := flag.String("snapshot", "", "Snapshot name to restore (required for restore)")
	source := flag.String("source", "", "Source directory (defaults to $HOME)")
	target := flag.String("target", "/var/backups/onevision", "Backup target directory")
	flag.Parse()

	if *source == "" {
		*source = os.Getenv("HOME")
	}

	config := backup.BackupConfig{
		SourceDir: *source,
		TargetDir: *target,
		Excluded:  []string{".cache", "Downloads", ".local/share/Trash"},
	}

	switch *action {
	case "backup":
		fmt.Printf("Starting backup of %s to %s...\n", config.SourceDir, config.TargetDir)
		path, err := backup.CreateSnapshot(config)
		if err != nil {
			log.Fatalf("Backup failed: %v\nOutput: %s", err, path)
		}
		fmt.Printf("Backup successful! Saved to: %s\n", path)

	case "restore":
		if *snapshot == "" {
			log.Fatal("Error: -snapshot name is required for restore action.")
		}
		snapshotPath := filepath.Join(config.TargetDir, *snapshot)
		fmt.Printf("Restoring from %s to %s...\n", snapshotPath, config.SourceDir)
		output, err := backup.RestoreSnapshot(snapshotPath, config.SourceDir)
		if err != nil {
			log.Fatalf("Restore failed: %v\nOutput: %s", err, output)
		}
		fmt.Println("Restore completed successfully!")

	case "list":
		snapshots, err := backup.ListSnapshots(config.TargetDir)
		if err != nil {
			log.Fatalf("Failed to list snapshots: %v", err)
		}
		fmt.Println("Available Snapshots:")
		for _, s := range snapshots {
			fmt.Printf(" - %s\n", s)
		}

	default:
		fmt.Printf("Unknown action: %s\n", *action)
		flag.Usage()
	}
}
