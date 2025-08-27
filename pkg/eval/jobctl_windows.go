//go:build windows

package eval

import (
	"fmt"
	"strconv"
	"sync"
	"syscall"
	"unsafe"

	"golang.org/x/sys/windows"
)

// Windows Job Objects API constants
const (
	// Job Object Limit Flags
	JOB_OBJECT_LIMIT_BREAKAWAY_OK      = 0x00000800
	JOB_OBJECT_LIMIT_KILL_ON_JOB_CLOSE = 0x00002000

	// Job Object Information Classes
	JobObjectBasicLimitInformation    = 2
	JobObjectExtendedLimitInformation = 9
)

// Windows API structures
type JOBOBJECT_BASIC_LIMIT_INFORMATION struct {
	PerProcessUserTimeLimit uint64
	PerJobUserTimeLimit     uint64
	LimitFlags              uint32
	MinimumWorkingSetSize   uintptr
	MaximumWorkingSetSize   uintptr
	ActiveProcessLimit      uint32
	Affinity                uintptr
	PriorityClass           uint32
	SchedulingClass         uint32
}

type JOBOBJECT_EXTENDED_LIMIT_INFORMATION struct {
	BasicLimitInformation JOBOBJECT_BASIC_LIMIT_INFORMATION
	IoInfo                struct {
		ReadOperationCount  uint64
		WriteOperationCount uint64
		OtherOperationCount uint64
		ReadTransferCount   uint64
		WriteTransferCount  uint64
		OtherTransferCount  uint64
	}
	ProcessMemoryLimit    uintptr
	JobMemoryLimit        uintptr
	PeakProcessMemoryUsed uintptr
	PeakJobMemoryUsed     uintptr
}

// Windows API functions
var (
	kernel32                      = windows.NewLazySystemDLL("kernel32.dll")
	procCreateJobObjectW          = kernel32.NewProc("CreateJobObjectW")
	procAssignProcessToJobObject  = kernel32.NewProc("AssignProcessToJobObject")
	procSetInformationJobObject   = kernel32.NewProc("SetInformationJobObject")
	procQueryInformationJobObject = kernel32.NewProc("QueryInformationJobObject")
	procTerminateJobObject        = kernel32.NewProc("TerminateJobObject")
)

// windowsJobController implements JobController for Windows using Job Objects.
type windowsJobController struct {
	mu   sync.RWMutex
	jobs map[windows.Handle]*windowsJob
}

// windowsJob represents a Windows Job Object.
type windowsJob struct {
	handle    windows.Handle
	processes map[uint32]bool // pid -> exists
	mu        sync.RWMutex
}

// windowsJobID implements JobID for Windows systems.
type windowsJobID struct {
	handle windows.Handle
}

func (j *windowsJobID) String() string {
	return fmt.Sprintf("job:0x%x", uintptr(j.handle))
}

func (j *windowsJobID) IsValid() bool {
	return j.handle != windows.InvalidHandle && j.handle != 0
}

// newPlatformJobController creates a Windows job controller.
func newPlatformJobController() (JobController, error) {
	return &windowsJobController{
		jobs: make(map[windows.Handle]*windowsJob),
	}, nil
}

func (c *windowsJobController) CreateJob() (JobID, error) {
	// Create a new job object
	ret, _, err := procCreateJobObjectW.Call(0, 0)
	if ret == 0 {
		return nil, err
	}

	handle := windows.Handle(ret)

	// Set job limits to allow breakaway and kill on job close
	var extendedInfo JOBOBJECT_EXTENDED_LIMIT_INFORMATION
	extendedInfo.BasicLimitInformation.LimitFlags = JOB_OBJECT_LIMIT_BREAKAWAY_OK | JOB_OBJECT_LIMIT_KILL_ON_JOB_CLOSE

	ret2, _, err2 := procSetInformationJobObject.Call(
		uintptr(handle),
		JobObjectExtendedLimitInformation,
		uintptr(unsafe.Pointer(&extendedInfo)),
		unsafe.Sizeof(extendedInfo),
	)
	if ret2 == 0 {
		windows.CloseHandle(handle)
		return nil, err2
	}

	// Track the job
	c.mu.Lock()
	c.jobs[handle] = &windowsJob{
		handle:    handle,
		processes: make(map[uint32]bool),
	}
	c.mu.Unlock()

	return &windowsJobID{handle: handle}, nil
}

func (c *windowsJobController) AddProcess(jobID JobID, pid int) error {
	winID, ok := jobID.(*windowsJobID)
	if !ok || !winID.IsValid() {
		return ErrInvalidJobID
	}

	// Open the process
	processHandle, err := windows.OpenProcess(
		windows.PROCESS_SET_QUOTA|windows.PROCESS_TERMINATE,
		false,
		uint32(pid),
	)
	if err != nil {
		return err
	}
	defer windows.CloseHandle(processHandle)

	// Assign process to job object
	ret, _, err := procAssignProcessToJobObject.Call(
		uintptr(winID.handle),
		uintptr(processHandle),
	)
	if ret == 0 {
		return err
	}

	// Track the process in our job
	c.mu.RLock()
	job, exists := c.jobs[winID.handle]
	c.mu.RUnlock()

	if !exists {
		return ErrJobNotFound
	}

	job.mu.Lock()
	job.processes[uint32(pid)] = true
	job.mu.Unlock()

	return nil
}

func (c *windowsJobController) BringToForeground(jobID JobID) error {
	// On Windows, we don't have the same concept of terminal foreground/background
	// as Unix. Instead, we can resume suspended processes and bring windows to front.
	// For now, we'll just resume the job processes.
	return c.ResumeJob(jobID)
}

func (c *windowsJobController) SendToBackground(jobID JobID) error {
	// On Windows, sending to "background" typically means minimizing windows
	// or reducing priority. For now, we'll just ensure processes continue running.
	// This is a no-op since Windows doesn't suspend processes by default.
	winID, ok := jobID.(*windowsJobID)
	if !ok || !winID.IsValid() {
		return ErrInvalidJobID
	}
	return nil
}

func (c *windowsJobController) ResumeJob(jobID JobID) error {
	winID, ok := jobID.(*windowsJobID)
	if !ok || !winID.IsValid() {
		return ErrInvalidJobID
	}

	// Get all processes in the job
	pids, err := c.GetJobProcesses(jobID)
	if err != nil {
		return err
	}

	// Resume each process (this is approximate - Windows doesn't have
	// exact equivalents to SIGCONT)
	for _, pid := range pids {
		processHandle, err := windows.OpenProcess(
			windows.PROCESS_SUSPEND_RESUME,
			false,
			uint32(pid),
		)
		if err != nil {
			continue // Skip processes we can't access
		}

		// Note: This is a simplified implementation. Full implementation would
		// need to track thread handles and use ResumeThread on each thread.
		windows.CloseHandle(processHandle)
	}

	return nil
}

func (c *windowsJobController) SuspendJob(jobID JobID) error {
	winID, ok := jobID.(*windowsJobID)
	if !ok || !winID.IsValid() {
		return ErrInvalidJobID
	}

	// Similar to ResumeJob, this would need to suspend all threads
	// in all processes in the job. This is a complex operation on Windows.

	// For now, return an error indicating this is not fully implemented
	return fmt.Errorf("SuspendJob not fully implemented on Windows")
}

func (c *windowsJobController) WaitForJob(jobID JobID) ([]ProcessStatus, error) {
	winID, ok := jobID.(*windowsJobID)
	if !ok || !winID.IsValid() {
		return nil, ErrInvalidJobID
	}

	// Get current processes in the job
	pids, err := c.GetJobProcesses(jobID)
	if err != nil {
		return nil, err
	}

	var statuses []ProcessStatus
	for _, pid := range pids {
		processHandle, err := windows.OpenProcess(
			windows.SYNCHRONIZE|windows.PROCESS_QUERY_INFORMATION,
			false,
			uint32(pid),
		)
		if err != nil {
			statuses = append(statuses, ProcessStatus{
				PID:   pid,
				Error: err,
			})
			continue
		}

		// Wait for the process to complete
		event, err := windows.WaitForSingleObject(processHandle, windows.INFINITE)
		if err != nil {
			windows.CloseHandle(processHandle)
			statuses = append(statuses, ProcessStatus{
				PID:   pid,
				Error: err,
			})
			continue
		}

		// Get exit code
		var exitCode uint32
		err = windows.GetExitCodeProcess(processHandle, &exitCode)
		windows.CloseHandle(processHandle)

		status := ProcessStatus{
			PID: pid,
		}
		if err != nil {
			status.Error = err
		} else {
			status.ExitCode = int(exitCode)
			status.Terminated = (event == windows.WAIT_OBJECT_0)
		}

		statuses = append(statuses, status)
	}

	return statuses, nil
}

func (c *windowsJobController) GetJobProcesses(jobID JobID) ([]int, error) {
	winID, ok := jobID.(*windowsJobID)
	if !ok || !winID.IsValid() {
		return nil, ErrInvalidJobID
	}

	c.mu.RLock()
	job, exists := c.jobs[winID.handle]
	c.mu.RUnlock()

	if !exists {
		return nil, ErrJobNotFound
	}

	job.mu.RLock()
	defer job.mu.RUnlock()

	pids := make([]int, 0, len(job.processes))
	for pid := range job.processes {
		pids = append(pids, int(pid))
	}

	return pids, nil
}

func (c *windowsJobController) TerminateJob(jobID JobID) error {
	winID, ok := jobID.(*windowsJobID)
	if !ok || !winID.IsValid() {
		return ErrInvalidJobID
	}

	// Terminate all processes in the job
	ret, _, err := procTerminateJobObject.Call(
		uintptr(winID.handle),
		1, // exit code
	)
	if ret == 0 {
		return err
	}

	// Clean up
	c.mu.Lock()
	delete(c.jobs, winID.handle)
	c.mu.Unlock()

	return windows.CloseHandle(winID.handle)
}

func (c *windowsJobController) Close() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	// Close all job handles
	for handle, job := range c.jobs {
		windows.CloseHandle(handle)
		job.mu.Lock()
		job.processes = make(map[uint32]bool)
		job.mu.Unlock()
	}

	c.jobs = make(map[windows.Handle]*windowsJob)
	return nil
}

// Helper function to create a job ID from a handle
func NewWindowsJobID(handle windows.Handle) JobID {
	return &windowsJobID{handle: handle}
}

// fgWindows implements the Windows-specific fg command using the job controller.
func fgWindows(controller JobController, pids ...int) error {
	if len(pids) == 0 {
		return ErrJobNotFound
	}

	// Create a new job for these processes
	jobID, err := controller.CreateJob()
	if err != nil {
		return err
	}

	// Add all processes to the job
	for _, pid := range pids {
		if err := controller.AddProcess(jobID, pid); err != nil {
			controller.TerminateJob(jobID) // Clean up on failure
			return err
		}
	}

	// Bring the job to foreground (resume processes)
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
			// Create a Windows equivalent of WaitStatus
			ws := syscall.WaitStatus{ExitCode: uint32(status.ExitCode)}
			errors[i] = &exception{NewExternalCmdExit(
				"[pid "+strconv.Itoa(status.PID)+"]",
				ws,
				status.PID), nil}
		}
	}

	return MakePipelineError(errors)
}
