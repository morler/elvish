//go:build unix

package eval

import (
	"fmt"
	"strconv"
	"sync"
	"syscall"

	"src.elv.sh/pkg/sys/eunix"
)

// unixJobController implements JobController for Unix systems using process groups.
type unixJobController struct {
	mu   sync.RWMutex
	jobs map[int]*unixJob // pgid -> job
}

// unixJob represents a Unix process group.
type unixJob struct {
	pgid      int
	processes map[int]bool // pid -> exists
	mu        sync.RWMutex
}

// unixJobID implements JobID for Unix systems.
type unixJobID struct {
	pgid int
}

func (j *unixJobID) String() string {
	return fmt.Sprintf("pgid:%d", j.pgid)
}

func (j *unixJobID) IsValid() bool {
	return j.pgid > 0
}

// newPlatformJobController creates a Unix job controller.
func newPlatformJobController() (JobController, error) {
	return &unixJobController{
		jobs: make(map[int]*unixJob),
	}, nil
}

func (c *unixJobController) CreateJob() (JobID, error) {
	// On Unix, we don't create process groups explicitly.
	// Instead, we track them when processes are added.
	// Return a placeholder that will be replaced when the first process is added.
	return &unixJobID{pgid: -1}, nil
}

func (c *unixJobController) AddProcess(jobID JobID, pid int) error {
	unixID, ok := jobID.(*unixJobID)
	if !ok {
		return ErrInvalidJobID
	}

	// Get the actual process group ID for this process
	pgid, err := syscall.Getpgid(pid)
	if err != nil {
		return err
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	// If this is a placeholder job ID, update it with the real pgid
	if unixID.pgid == -1 {
		unixID.pgid = pgid
	} else if unixID.pgid != pgid {
		// Process is not in the expected process group
		return ErrProcessNotInJob
	}

	// Find or create the job
	job, exists := c.jobs[pgid]
	if !exists {
		job = &unixJob{
			pgid:      pgid,
			processes: make(map[int]bool),
		}
		c.jobs[pgid] = job
	}

	job.mu.Lock()
	job.processes[pid] = true
	job.mu.Unlock()

	return nil
}

func (c *unixJobController) BringToForeground(jobID JobID) error {
	unixID, ok := jobID.(*unixJobID)
	if !ok || !unixID.IsValid() {
		return ErrInvalidJobID
	}

	pgid := unixID.pgid

	// Set the process group as the foreground process group
	if err := eunix.Tcsetpgrp(0, pgid); err != nil {
		return err
	}

	// Send SIGCONT to resume all processes in the group
	if err := syscall.Kill(-pgid, syscall.SIGCONT); err != nil {
		return err
	}

	return nil
}

func (c *unixJobController) SendToBackground(jobID JobID) error {
	unixID, ok := jobID.(*unixJobID)
	if !ok || !unixID.IsValid() {
		return ErrInvalidJobID
	}

	// On Unix, sending to background means giving terminal control back to shell
	// and ensuring the process group continues running
	shellPgid := syscall.Getpgrp()
	if err := eunix.Tcsetpgrp(0, shellPgid); err != nil {
		return err
	}

	// Ensure the background job continues running
	return syscall.Kill(-unixID.pgid, syscall.SIGCONT)
}

func (c *unixJobController) ResumeJob(jobID JobID) error {
	unixID, ok := jobID.(*unixJobID)
	if !ok || !unixID.IsValid() {
		return ErrInvalidJobID
	}

	return syscall.Kill(-unixID.pgid, syscall.SIGCONT)
}

func (c *unixJobController) SuspendJob(jobID JobID) error {
	unixID, ok := jobID.(*unixJobID)
	if !ok || !unixID.IsValid() {
		return ErrInvalidJobID
	}

	return syscall.Kill(-unixID.pgid, syscall.SIGSTOP)
}

func (c *unixJobController) WaitForJob(jobID JobID) ([]ProcessStatus, error) {
	unixID, ok := jobID.(*unixJobID)
	if !ok || !unixID.IsValid() {
		return nil, ErrInvalidJobID
	}

	c.mu.RLock()
	job, exists := c.jobs[unixID.pgid]
	c.mu.RUnlock()

	if !exists {
		return nil, ErrJobNotFound
	}

	job.mu.RLock()
	pids := make([]int, 0, len(job.processes))
	for pid := range job.processes {
		pids = append(pids, pid)
	}
	job.mu.RUnlock()

	var statuses []ProcessStatus
	for _, pid := range pids {
		var ws syscall.WaitStatus
		wpid, err := syscall.Wait4(pid, &ws, syscall.WUNTRACED, nil)

		status := ProcessStatus{
			PID: wpid,
		}

		if err != nil {
			status.Error = err
		} else if ws.Exited() {
			status.ExitCode = ws.ExitStatus()
		} else if ws.Signaled() {
			status.Terminated = true
			status.ExitCode = int(ws.Signal())
		}

		statuses = append(statuses, status)
	}

	return statuses, nil
}

func (c *unixJobController) GetJobProcesses(jobID JobID) ([]int, error) {
	unixID, ok := jobID.(*unixJobID)
	if !ok || !unixID.IsValid() {
		return nil, ErrInvalidJobID
	}

	c.mu.RLock()
	job, exists := c.jobs[unixID.pgid]
	c.mu.RUnlock()

	if !exists {
		return nil, ErrJobNotFound
	}

	job.mu.RLock()
	defer job.mu.RUnlock()

	pids := make([]int, 0, len(job.processes))
	for pid := range job.processes {
		pids = append(pids, pid)
	}

	return pids, nil
}

func (c *unixJobController) TerminateJob(jobID JobID) error {
	unixID, ok := jobID.(*unixJobID)
	if !ok || !unixID.IsValid() {
		return ErrInvalidJobID
	}

	// Send SIGTERM to the entire process group
	if err := syscall.Kill(-unixID.pgid, syscall.SIGTERM); err != nil {
		return err
	}

	// Clean up from our tracking
	c.mu.Lock()
	delete(c.jobs, unixID.pgid)
	c.mu.Unlock()

	return nil
}

func (c *unixJobController) Close() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	// Clear all tracked jobs
	c.jobs = make(map[int]*unixJob)
	return nil
}

// Helper function to create a job ID from a process group ID
func NewUnixJobID(pgid int) JobID {
	return &unixJobID{pgid: pgid}
}

// fgUnix implements the Unix-specific fg command using the job controller.
func fgUnix(controller JobController, pids ...int) error {
	if len(pids) == 0 {
		return ErrJobNotFound
	}

	// Validate all PIDs are in the same process group
	var pgid int
	for i, pid := range pids {
		currentPgid, err := syscall.Getpgid(pid)
		if err != nil {
			return err
		}
		if i == 0 {
			pgid = currentPgid
		} else if currentPgid != pgid {
			return ErrProcessNotInJob
		}
	}

	// Create a job ID and add all processes
	jobID := NewUnixJobID(pgid)
	for _, pid := range pids {
		if err := controller.AddProcess(jobID, pid); err != nil {
			return err
		}
	}

	// Bring the job to foreground
	if err := controller.BringToForeground(jobID); err != nil {
		return err
	}

	// Wait for the job to complete
	statuses, err := controller.WaitForJob(jobID)
	if err != nil {
		return err
	}

	// Convert statuses to exceptions (matching original fg behavior)
	errors := make([]Exception, len(statuses))
	for i, status := range statuses {
		if status.Error != nil {
			errors[i] = &exception{status.Error, nil}
		} else {
			errors[i] = &exception{NewExternalCmdExit(
				"[pid "+strconv.Itoa(status.PID)+"]",
				syscall.WaitStatus(status.ExitCode),
				status.PID), nil}
		}
	}

	return MakePipelineError(errors)
}
