//go:build unix

package term

import (
	"time"
)

type byteReaderWithTimeout interface {
	// ReadByteWithTimeout reads a single byte with a timeout. A negative
	// timeout means no timeout.
	ReadByteWithTimeout(timeout time.Duration) (byte, error)
}

// badRune is now defined in utf8_decode.go for cross-platform use

var utf8SeqTimeout = 10 * time.Millisecond

// Reads a rune from the reader. The timeout applies to the first byte; a
// negative value means no timeout.
func readRune(rd byteReaderWithTimeout, timeout time.Duration) (rune, error) {
	// Read the leader byte with the specified timeout
	leader, err := rd.ReadByteWithTimeout(timeout)
	if err != nil {
		return badRune, err
	}
	
	// Determine how many bytes we need for this UTF-8 sequence
	var needed int
	switch {
	case leader>>7 == 0:
		needed = 1
	case leader>>5 == 0x6:
		needed = 2
	case leader>>4 == 0xe:
		needed = 3
	case leader>>3 == 0x1e:
		needed = 4
	default:
		return badRune, errInvalidUTF8
	}
	
	// Collect all bytes for the UTF-8 sequence
	bytes := make([]byte, needed)
	bytes[0] = leader
	
	// Read remaining bytes with UTF-8 sequence timeout
	for i := 1; i < needed; i++ {
		b, err := rd.ReadByteWithTimeout(utf8SeqTimeout)
		if err != nil {
			return badRune, err
		}
		bytes[i] = b
	}
	
	// Use the cross-platform UTF-8 decoding logic
	r, _, err := decodeUTF8FromBytes(bytes)
	return r, err
}
