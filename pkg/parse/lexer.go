package parse

import "unicode/utf8"

// EOF is the sentinel rune value used to indicate end of input.
const EOF rune = -1

// Lexer provides basic lexing functionality for parsing text.
// It can be embedded in parser structures to provide common text reading operations.
type Lexer struct {
	src     string // source text being parsed
	pos     int    // current reading position
	overEOF int    // count of reads past EOF
}

// NewLexer creates a new Lexer for the given source text.
func NewLexer(src string) *Lexer {
	return &Lexer{src: src}
}

// Next reads and returns the next rune from the source.
// Returns EOF when the end of input is reached.
func (l *Lexer) Next() rune {
	if l.pos == len(l.src) {
		l.overEOF++
		return EOF
	}
	r, s := utf8.DecodeRuneInString(l.src[l.pos:])
	l.pos += s
	return r
}

// Backup moves the reading position back by one rune.
// Properly handles the EOF state.
func (l *Lexer) Backup() {
	if l.overEOF > 0 {
		l.overEOF--
		return
	}
	_, s := utf8.DecodeLastRuneInString(l.src[:l.pos])
	l.pos -= s
}

// Pos returns the current reading position.
func (l *Lexer) Pos() int {
	return l.pos
}

// Src returns the source string being parsed.
func (l *Lexer) Src() string {
	return l.src
}

// Reset resets the lexer state to the beginning of the source.
func (l *Lexer) Reset() {
	l.pos = 0
	l.overEOF = 0
}
