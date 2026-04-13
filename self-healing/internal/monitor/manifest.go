package monitor

import (
	"log"
)

type CriticalFile struct {
	Path         string
	ExpectedHash string
	Description  string
}

// Manifest defines the set of files that must be monitored for integrity
// In a production environment, this would be loaded from a signed config file
var Manifest = []CriticalFile{
	{
		Path:        "config.json",
		Description: "Simulated critical config file",
	},
	{
		Path:        "auth.key",
		Description: "Simulated security key",
	},
}

// InitializeManifest calculates initial hashes for the manifest files
// This assumes the files are currently in a "good" state.
func InitializeManifest() {
	for i, file := range Manifest {
		hash, err := Checksum(file.Path)
		if err != nil {
			log.Printf("[Monitor] Warning: Could not hash critical file %s: %v", file.Path, err)
			continue
		}
		Manifest[i].ExpectedHash = hash
		log.Printf("[Monitor] Manifest initialized: %s (%s)", file.Path, hash[:8])
	}
}
