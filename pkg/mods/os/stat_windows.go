package os

import (
	"syscall"
	"time"

	"src.elv.sh/pkg/eval/vals"
)

// Taken from
// https://learn.microsoft.com/en-us/windows/win32/fileio/file-attribute-constants.
// The syscall package only has a subset of these.
//
// Some of these attributes are redundant with fields in the outer stat map, but
// we keep all of them for consistency.
var fileAttributeNames = [...]struct {
	bit  uint32
	name string
}{
	{0x1, "readonly"},
	{0x2, "hidden"},
	{0x4, "system"},
	{0x10, "directory"},
	{0x20, "archive"},
	{0x40, "device"},
	{0x80, "normal"},
	{0x100, "temporary"},
	{0x200, "sparse-file"},
	{0x400, "reparse-point"},
	{0x800, "compressed"},
	{0x1000, "offline"},
	{0x2000, "not-content-indexed"},
	{0x4000, "encrypted"},
	{0x8000, "integrity-system"},
	{0x10000, "virtual"},
	{0x20000, "no-scrub-data"},
	{0x40000, "ea"},
	{0x80000, "pinned"},
	{0x100000, "unpinned"},
	{0x400000, "recall-on-data-access"},
}

// filetimeToTime converts a Windows FILETIME to a Go time.Time.
// FILETIME represents the number of 100-nanosecond intervals since January 1, 1601 UTC.
func filetimeToTime(ft syscall.Filetime) time.Time {
	// Windows epoch starts at January 1, 1601 UTC
	// Unix epoch starts at January 1, 1970 UTC
	// The difference is 116444736000000000 * 100ns = 11644473600 seconds
	const windowsEpochDiff = 116444736000000000

	// Combine low and high parts into a single int64
	nsec := int64(ft.HighDateTime)<<32 + int64(ft.LowDateTime)

	// Convert to Unix epoch (nanoseconds since January 1, 1970 UTC)
	nsec = (nsec - windowsEpochDiff) * 100

	return time.Unix(0, nsec).UTC()
}

func statSysMap(sys any) vals.Map {
	attrData := sys.(*syscall.Win32FileAttributeData)
	// TODO: Make this a set when Elvish has a set type.
	fileAttributes := vals.EmptyList
	for _, attr := range fileAttributeNames {
		if attrData.FileAttributes&attr.bit != 0 {
			fileAttributes = fileAttributes.Conj(attr.name)
		}
	}
	return vals.MakeMap(
		"file-attributes", fileAttributes,
		"creation-time", filetimeToTime(attrData.CreationTime),
		"last-access-time", filetimeToTime(attrData.LastAccessTime),
		"last-write-time", filetimeToTime(attrData.LastWriteTime),
	)
}
