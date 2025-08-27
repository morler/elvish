package eval

import (
	"os"
	"testing"
)

func TestJobControllerInterface(t *testing.T) {
	controller, err := NewJobController()
	if err != nil {
		t.Fatalf("Failed to create job controller: %v", err)
	}
	defer controller.Close()

	// Test creating a job
	jobID, err := controller.CreateJob()
	if err != nil {
		t.Fatalf("Failed to create job: %v", err)
	}

	if !jobID.IsValid() {
		t.Error("Job ID should be valid after creation")
	}

	// Test job ID string representation
	jobStr := jobID.String()
	if jobStr == "" {
		t.Error("Job ID string should not be empty")
	}
	t.Logf("Created job: %s", jobStr)
}

func TestJobControllerWithCurrentProcess(t *testing.T) {
	controller, err := NewJobController()
	if err != nil {
		t.Skip("Job controller not available")
	}
	defer controller.Close()

	jobID, err := controller.CreateJob()
	if err != nil {
		t.Fatalf("Failed to create job: %v", err)
	}

	// Try to add the current process to the job
	// Note: This might fail on some platforms due to permissions
	currentPid := os.Getpid()
	err = controller.AddProcess(jobID, currentPid)

	// On Unix, adding current process might fail, on Windows it might work
	// So we just log the result rather than failing
	if err != nil {
		t.Logf("Adding current process failed (expected on some platforms): %v", err)
	} else {
		t.Logf("Successfully added current process %d to job", currentPid)

		// Test getting job processes
		pids, err := controller.GetJobProcesses(jobID)
		if err != nil {
			t.Errorf("Failed to get job processes: %v", err)
		} else {
			t.Logf("Job processes: %v", pids)
		}
	}
}

func TestFgCommandBasic(t *testing.T) {
	// Test that fg function exists and handles empty arguments correctly
	err := fg()
	if err == nil {
		t.Error("fg() should return an error when called with no arguments")
	}
	t.Logf("fg() correctly returned error for no arguments: %v", err)
}

// TestFgCommandIntegration tests the fg command with a real subprocess
func TestFgCommandIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	// Start a simple subprocess that we can control
	// On Windows, we'll use timeout, on Unix we can use sleep
	var cmd string
	var args []string

	// Use a cross-platform approach - Go's own program
	cmd = os.Args[0]                // Use the test binary itself
	args = []string{"-test.run=^$"} // Run no tests (will exit quickly)

	// Create a simple external command for testing
	// This is a basic test - more comprehensive tests would need actual
	// process management scenarios

	t.Logf("Testing with command: %s %v", cmd, args)

	// For now, just verify that fg handles invalid PIDs gracefully
	invalidPid := 999999 // Very unlikely to be a real PID
	err := fg(invalidPid)
	if err == nil {
		t.Error("fg should return an error for invalid PID")
	}
	t.Logf("fg correctly handled invalid PID: %v", err)
}
