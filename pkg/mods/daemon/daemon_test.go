package daemon_test

import (
	"embed"
	"testing"

	"src.elv.sh/pkg/daemon/daemondefs"
	"src.elv.sh/pkg/eval"
	"src.elv.sh/pkg/eval/evaltest"
	"src.elv.sh/pkg/mods/daemon"
	"src.elv.sh/pkg/store/storedefs"
)

//go:embed *.elvts
var transcripts embed.FS

func TestTranscripts(t *testing.T) {
	evaltest.TestTranscriptsInFS(t, transcripts,
		"use-daemon", func(t *testing.T, ev *eval.Evaler) {
			ev.ExtendGlobal(eval.BuildNs().AddNs("daemon", daemon.Ns(newMockClient())))
		},
	)
}

func TestDaemon(t *testing.T) {
	client := newMockClient()
	ns := daemon.Ns(client)

	// Test pid variable
	pidVar, ok := ns.Index("pid")
	if !ok {
		t.Error("daemon module should export pid variable")
	}
	if pidVar == nil {
		t.Error("pid variable should not be nil")
	}

	// Test sock variable
	sockVar, ok := ns.Index("sock")
	if !ok {
		t.Error("daemon module should export sock variable")
	}
	if sockVar == nil {
		t.Error("sock variable should not be nil")
	}
}

// mockClient implements daemondefs.Client for testing
type mockClient struct{}

func newMockClient() daemondefs.Client {
	return &mockClient{}
}

func (c *mockClient) Pid() (int, error) {
	return 12345, nil
}

func (c *mockClient) SockPath() string {
	return "/tmp/elvish-test.sock"
}

func (c *mockClient) Version() (int, error) {
	return 1, nil
}

func (c *mockClient) ResetConn() error {
	return nil
}

func (c *mockClient) Close() error {
	return nil
}

// Store interface methods (required by daemondefs.Client which embeds storedefs.Store)
func (c *mockClient) NextCmdSeq() (int, error) {
	return 1, nil
}

func (c *mockClient) AddCmd(text string) (int, error) {
	return 1, nil
}

func (c *mockClient) DelCmd(seq int) error {
	return nil
}

func (c *mockClient) Cmd(seq int) (string, error) {
	return "echo test", nil
}

func (c *mockClient) CmdsWithSeq(from, upto int) ([]storedefs.Cmd, error) {
	return []storedefs.Cmd{
		{Text: "echo test", Seq: from},
	}, nil
}

func (c *mockClient) NextCmd(from int, prefix string) (storedefs.Cmd, error) {
	return storedefs.Cmd{Text: prefix + "test", Seq: from + 1}, nil
}

func (c *mockClient) PrevCmd(upto int, prefix string) (storedefs.Cmd, error) {
	return storedefs.Cmd{Text: prefix + "test", Seq: upto - 1}, nil
}

func (c *mockClient) AddDir(dir string, incFactor float64) error {
	return nil
}

func (c *mockClient) DelDir(dir string) error {
	return nil
}

func (c *mockClient) Dirs(blacklist map[string]struct{}) ([]storedefs.Dir, error) {
	return []storedefs.Dir{
		{Path: "/home/test", Score: 10.0},
	}, nil
}