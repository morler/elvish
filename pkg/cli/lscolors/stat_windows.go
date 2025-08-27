package lscolors

import (
	"os"
)

// Note: This implementation is conservative and returns false.
// Windows does support hard links, but detecting them requires additional
// system calls that may impact performance. The hard link count is not
// available directly from os.FileInfo on Windows like it is on Unix systems.
//
// A complete implementation would need:
// 1. Opening the file with CreateFile
// 2. Calling GetFileInformationByHandle
// 3. Checking the NumberOfLinks field
//
// For now, we prioritize performance and correctness over feature completeness.
func isMultiHardlink(info os.FileInfo) bool {
	// Windows supports hardlinks, but detecting them requires additional I/O.
	// For performance reasons, we return false. This means hard-linked files
	// won't get special highlighting in ls colors, but functionality remains correct.
	return false
}
