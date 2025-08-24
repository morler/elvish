package term

// Cross-platform UTF-8 decoding logic extracted from readRune
// This can be tested without platform-specific dependencies

// Moved from read_rune.go to make it cross-platform
const badRune = '\ufffd'

// decodeUTF8FromBytes decodes a single UTF-8 rune from a byte slice.
// Returns the decoded rune, number of bytes consumed, and any error.
// This function is cross-platform and doesn't require Unix-specific functionality.
func decodeUTF8FromBytes(data []byte) (rune, int, error) {
	if len(data) == 0 {
		return badRune, 0, errEOF
	}

	leader := data[0]
	var r rune
	var needed int

	switch {
	case leader>>7 == 0:
		// ASCII character (0xxxxxxx)
		r = rune(leader)
		needed = 1
	case leader>>5 == 0x6:
		// 2-byte UTF-8 (110xxxxx)
		r = rune(leader & 0x1f)
		needed = 2
	case leader>>4 == 0xe:
		// 3-byte UTF-8 (1110xxxx)
		r = rune(leader & 0xf)
		needed = 3
	case leader>>3 == 0x1e:
		// 4-byte UTF-8 (11110xxx)
		r = rune(leader & 0x7)
		needed = 4
	default:
		// Invalid UTF-8 leading byte
		return badRune, 1, errInvalidUTF8
	}

	// Check if we have enough bytes, but validate continuation bytes first
	// to distinguish between incomplete and invalid sequences
	for i := 1; i < needed && i < len(data); i++ {
		b := data[i]
		if b>>6 != 0x2 { // Must be 10xxxxxx
			return badRune, i, errInvalidUTF8
		}
	}
	
	if len(data) < needed {
		return badRune, len(data), errIncompleteUTF8
	}

	// Read continuation bytes (we've already validated them above)
	for i := 1; i < needed; i++ {
		b := data[i]
		r = r<<6 + rune(b&0x3f)
	}

	return r, needed, nil
}

// Common errors for UTF-8 decoding
var (
	errEOF           = &utf8Error{"unexpected end of data"}
	errInvalidUTF8   = &utf8Error{"invalid UTF-8 sequence"}
	errIncompleteUTF8 = &utf8Error{"incomplete UTF-8 sequence"}
)

type utf8Error struct {
	msg string
}

func (e *utf8Error) Error() string {
	return e.msg
}