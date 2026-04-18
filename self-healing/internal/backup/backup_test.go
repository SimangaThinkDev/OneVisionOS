package backup

import (
	"os"
	"path/filepath"
	"testing"
)

func TestCreateAndListSnapshots(t *testing.T) {
	// Setup temporary directories
	tempDir, err := os.MkdirTemp("", "onevision-backup-test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tempDir)

	sourceDir := filepath.Join(tempDir, "source")
	targetDir := filepath.Join(tempDir, "target")
	os.Mkdir(sourceDir, 0755)
	os.Mkdir(targetDir, 0755)

	// Create a dummy file in source
	testFile := filepath.Join(sourceDir, "test.txt")
	os.WriteFile(testFile, []byte("hello backup"), 0644)

	config := BackupConfig{
		SourceDir: sourceDir,
		TargetDir: targetDir,
		Excluded:  []string{},
	}

	// 1. Create Snapshot
	snapshotPath, err := CreateSnapshot(config)
	if err != nil {
		t.Fatalf("CreateSnapshot failed: %v", err)
	}

	if _, err := os.Stat(snapshotPath); os.IsNotExist(err) {
		t.Errorf("Snapshot directory was not created: %s", snapshotPath)
	}

	// Check if file exists in snapshot
	snapFile := filepath.Join(snapshotPath, "test.txt")
	if _, err := os.Stat(snapFile); os.IsNotExist(err) {
		t.Errorf("File not found in snapshot: %s", snapFile)
	}

	// 2. List Snapshots
	snapshots, err := ListSnapshots(targetDir)
	if err != nil {
		t.Fatalf("ListSnapshots failed: %v", err)
	}

	if len(snapshots) != 1 {
		t.Errorf("Expected 1 snapshot, got %d", len(snapshots))
	}

	// 3. Test 'latest' symlink
	latest := filepath.Join(targetDir, "latest")
	if _, err := os.Lstat(latest); err != nil {
		t.Errorf("'latest' symlink not created: %v", err)
	}
}

func TestRestoreSnapshot(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "onevision-restore-test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tempDir)

	snapDir := filepath.Join(tempDir, "snap")
	restoreDir := filepath.Join(tempDir, "restore")
	os.Mkdir(snapDir, 0755)
	os.Mkdir(restoreDir, 0755)

	os.WriteFile(filepath.Join(snapDir, "recovered.txt"), []byte("i am back"), 0644)

	_, err = RestoreSnapshot(snapDir, restoreDir)
	if err != nil {
		t.Fatalf("RestoreSnapshot failed: %v", err)
	}

	recoveredFile := filepath.Join(restoreDir, "recovered.txt")
	content, err := os.ReadFile(recoveredFile)
	if err != nil {
		t.Fatalf("Failed to read recovered file: %v", err)
	}

	if string(content) != "i am back" {
		t.Errorf("Recovered content mismatch: got %s", string(content))
	}
}
