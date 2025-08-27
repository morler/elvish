package os_test

import (
	"embed"
	"net"
	"os"
	"strconv"
	"testing"

	"src.elv.sh/pkg/eval/evaltest"
	"src.elv.sh/pkg/must"
	"src.elv.sh/pkg/testutil"
)

//go:embed *.elvts
var transcripts embed.FS

func TestTranscripts(t *testing.T) {
	evaltest.TestTranscriptsInFS(t, transcripts,
		"umask", func(t *testing.T, arg string) {
			testutil.Umask(t, must.OK1(strconv.Atoi(arg)))
		},
		"mkfifo-or-skip", mkFifoOrSkip,
		"mksock-or-skip", func(t *testing.T, s string) {
			listener, err := net.Listen("unix", "./sock")
			if err != nil {
				t.Skipf("can't listen to UNIX socket: %v", err)
			}
			// Ensure socket file is accessible and is indeed a socket
			info, err := os.Stat("./sock")
			if err != nil {
				t.Skipf("socket file not accessible: %v", err)
			}
			if info.Mode()&os.ModeSocket == 0 {
				t.Skipf("socket file is not a socket (mode: %v)", info.Mode())
			}
			t.Cleanup(func() {
				listener.Close()
				// On Windows, socket file is auto-deleted when listener closes
				// On Unix, we might need to manually remove it
				os.Remove("./sock")
			})
		},
		"only-if-can-create-symlink", func(t *testing.T) {
			testutil.ApplyDir(testutil.Dir{"test-file": ""})
			err := os.Symlink("test-file", "test-symlink")
			if err != nil {
				// On Windows we may or may not be able to create a symlink.
				t.Skipf("symlink: %v", err)
			}
			must.OK(os.Remove("test-file"))
			must.OK(os.Remove("test-symlink"))
		},
		"create-windows-special-files-or-skip", createWindowsSpecialFileOrSkip,
	)
}
