package eval

import "syscall"

// Nop on Windows.
func putSelfInFg() error { return nil }

// Windows CreateProcess flags for proper process isolation
const (
	detachedProcess        = 0x00000008 // Start process in background
	createBreakawayFromJob = 0x01000000 // Break away from job object (for better cleanup)
	createNewProcessGroup  = 0x00000200 // Create new process group
)

func makeSysProcAttr(bg bool) *syscall.SysProcAttr {
	// Always use flags that help with process cleanup on Windows
	flags := uint32(createBreakawayFromJob | createNewProcessGroup)

	if bg {
		flags |= detachedProcess
	}

	return &syscall.SysProcAttr{CreationFlags: flags}
}
