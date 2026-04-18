package update

import (
	"os"
	"os/exec"
	"testing"
)

// Helper to mock exec.Command
func mockCommand(command string, args ...string) *exec.Cmd {
	cs := []string{"-test.run=TestHelperProcess", "--", command}
	cs = append(cs, args...)
	cmd := exec.Command(os.Args[0], cs...)
	cmd.Env = []string{"GO_WANT_HELPER_PROCESS=1"}
	return cmd
}

func TestHelperProcess(t *testing.T) {
	if os.Getenv("GO_WANT_HELPER_PROCESS") != "1" {
		return
	}
	// Successful Exit
	os.Exit(0)
}

func TestRunUpdate(t *testing.T) {
	// Mock the runner
	oldRunner := commandRunner
	commandRunner = mockCommand
	defer func() { commandRunner = oldRunner }()

	err := RunUpdate()
	if err != nil {
		t.Errorf("RunUpdate failed: %v", err)
	}
}

func TestCheckForUpdates(t *testing.T) {
	oldRunner := commandRunner
	commandRunner = mockCommand
	defer func() { commandRunner = oldRunner }()

	hasUpdates, err := CheckForUpdates()
	if err != nil {
		t.Errorf("CheckForUpdates failed: %v", err)
	}
	if !hasUpdates {
		t.Errorf("Expected updates available")
	}
}
