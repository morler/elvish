package eval

import (
	"errors"
)

// JobController provides cross-platform job control functionality.
// On Unix, it manages process groups using signals and terminal control.
// On Windows, it uses Job Objects to manage groups of processes.
type JobController interface {
	// CreateJob creates a new job/process group and returns a job ID.
	CreateJob() (JobID, error)
	
	// AddProcess adds a process to an existing job/process group.
	AddProcess(jobID JobID, pid int) error
	
	// BringToForeground brings a job/process group to the foreground.
	// This is the core functionality needed by the fg command.
	BringToForeground(jobID JobID) error
	
	// SendToBackground sends a job/process group to the background.
	SendToBackground(jobID JobID) error
	
	// ResumeJob resumes all processes in a suspended job/process group.
	ResumeJob(jobID JobID) error
	
	// SuspendJob suspends all processes in a job/process group.
	SuspendJob(jobID JobID) error
	
	// WaitForJob waits for all processes in a job to complete.
	WaitForJob(jobID JobID) ([]ProcessStatus, error)
	
	// GetJobProcesses returns the list of process IDs in a job.
	GetJobProcesses(jobID JobID) ([]int, error)
	
	// TerminateJob forcefully terminates all processes in a job.
	TerminateJob(jobID JobID) error
	
	// Close cleans up resources associated with the job controller.
	Close() error
}

// JobID uniquely identifies a job/process group across platforms.
// On Unix, this corresponds to a process group ID (pgid).
// On Windows, this corresponds to a Job Object handle.
type JobID interface {
	// String returns a human-readable representation of the job ID.
	String() string
	
	// IsValid returns true if the job ID is valid.
	IsValid() bool
}

// ProcessStatus represents the status of a process after it completes.
type ProcessStatus struct {
	PID        int
	ExitCode   int
	Terminated bool // true if terminated by signal/force
	Error      error
}

// Common errors for job control operations.
var (
	ErrJobNotFound         = errors.New("job not found")
	ErrProcessNotInJob     = errors.New("process not in job")
	ErrJobControlNotSupported = errors.New("job control not supported on this platform")
	ErrInvalidJobID        = errors.New("invalid job ID")
	ErrJobAlreadyTerminated = errors.New("job already terminated")
)

// NewJobController creates a new job controller for the current platform.
func NewJobController() (JobController, error) {
	return newPlatformJobController()
}