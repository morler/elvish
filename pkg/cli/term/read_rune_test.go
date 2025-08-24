//go:build unix

package term

import "testing"

// These tests are Unix-specific because readRune is Unix-specific.
// The UTF-8 decoding logic is now tested cross-platform in utf8_decode_test.go.

var contents = []string{
	"English",
	"Ελληνικά",
	"你好 こんにちは",
	"𐌰𐌱",
}

func TestReadRune(t *testing.T) {
	for _, content := range contents {
		t.Run(content, func(t *testing.T) {
			rd, w, cleanup := setupFileReader()
			defer cleanup()

			w.Write([]byte(content))
			for _, wantRune := range content {
				r, err := readRune(rd, 0)
				if r != wantRune {
					t.Errorf("got rune %q, want %q", r, wantRune)
				}
				if err != nil {
					t.Errorf("got err %v, want nil", err)
				}
			}
		})
	}
}

func TestReadRune_ErrorAtFirstByte(t *testing.T) {
	rd, _, cleanup := setupFileReader()
	defer cleanup()

	r, err := readRune(rd, 0)
	if r != '\ufffd' {
		t.Errorf("got rune %q, want %q", r, '\ufffd')
	}
	if err == nil {
		t.Errorf("got err %v, want non-nil", err)
	}
}

func TestReadRune_ErrorAtNonFirstByte(t *testing.T) {
	rd, w, cleanup := setupFileReader()
	defer cleanup()

	w.Write([]byte{0xe4})

	r, err := readRune(rd, 0)
	if r != '\ufffd' {
		t.Errorf("got rune %q, want %q", r, '\ufffd')
	}
	if err == nil {
		t.Errorf("got err %v, want non-nil", err)
	}
}
