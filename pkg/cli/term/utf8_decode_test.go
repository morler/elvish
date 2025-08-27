package term

import (
	"testing"
)

var testContents = []string{
	"English",
	"Ελληνικά",
	"你好 こんにちは",
	"𐌰𐌱",
}

// Cross-platform test for UTF-8 decoding logic
func TestDecodeUTF8FromBytes(t *testing.T) {
	for _, content := range testContents {
		t.Run(content, func(t *testing.T) {
			data := []byte(content)
			offset := 0

			for _, expectedRune := range content {
				if offset >= len(data) {
					t.Errorf("unexpected end of data at offset %d", offset)
					break
				}

				r, consumed, err := decodeUTF8FromBytes(data[offset:])
				if err != nil {
					t.Errorf("decoding error at offset %d: %v", offset, err)
					break
				}
				if r != expectedRune {
					t.Errorf("got rune %q, want %q at offset %d", r, expectedRune, offset)
				}
				if consumed <= 0 {
					t.Errorf("consumed %d bytes, expected > 0", consumed)
					break
				}

				offset += consumed
			}

			if offset != len(data) {
				t.Errorf("consumed %d bytes, expected %d", offset, len(data))
			}
		})
	}
}

func TestDecodeUTF8FromBytes_EmptyData(t *testing.T) {
	r, consumed, err := decodeUTF8FromBytes([]byte{})
	if r != badRune {
		t.Errorf("got rune %q, want %q", r, badRune)
	}
	if consumed != 0 {
		t.Errorf("consumed %d bytes, want 0", consumed)
	}
	if err != errEOF {
		t.Errorf("got err %v, want %v", err, errEOF)
	}
}

func TestDecodeUTF8FromBytes_IncompleteSequence(t *testing.T) {
	// 0xe4 is the start of a 3-byte sequence but we only provide 1 byte
	r, consumed, err := decodeUTF8FromBytes([]byte{0xe4})
	if r != badRune {
		t.Errorf("got rune %q, want %q", r, badRune)
	}
	if consumed != 1 {
		t.Errorf("consumed %d bytes, want 1", consumed)
	}
	if err != errIncompleteUTF8 {
		t.Errorf("got err %v, want %v", err, errIncompleteUTF8)
	}
}

func TestDecodeUTF8FromBytes_InvalidSequence(t *testing.T) {
	// 0xff is not a valid UTF-8 leading byte
	r, consumed, err := decodeUTF8FromBytes([]byte{0xff})
	if r != badRune {
		t.Errorf("got rune %q, want %q", r, badRune)
	}
	if consumed != 1 {
		t.Errorf("consumed %d bytes, want 1", consumed)
	}
	if err != errInvalidUTF8 {
		t.Errorf("got err %v, want %v", err, errInvalidUTF8)
	}
}

func TestDecodeUTF8FromBytes_InvalidContinuation(t *testing.T) {
	// 0xe4 followed by 0xff (invalid continuation byte)
	r, consumed, err := decodeUTF8FromBytes([]byte{0xe4, 0xff})
	if r != badRune {
		t.Errorf("got rune %q, want %q", r, badRune)
	}
	if consumed != 1 {
		t.Errorf("consumed %d bytes, want 1", consumed)
	}
	if err != errInvalidUTF8 {
		t.Errorf("got err %v, want %v", err, errInvalidUTF8)
	}
}
