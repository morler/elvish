package parse

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"strings"
	"unicode/utf8"

	"src.elv.sh/pkg/diag"
)

// parser maintains some mutable states of parsing.
//
// NOTE: The src member is assumed to be valid UF-8.
type parser struct {
	*Lexer              // embedded lexer for basic text reading
	srcName string      // source name for error reporting
	errors  []*Error   // accumulated parse errors
	warn    io.Writer  // warning output writer
}

// Error is a parse error.
type Error = diag.Error[ErrorTag]

// ErrorTag parameterizes [diag.Error] to define [Error].
type ErrorTag struct{}

func (ErrorTag) ErrorTag() string { return "parse error" }

func parse[N Node](ps *parser, n N) parsed[N] {
	begin := ps.Pos()
	n.n().From = begin
	n.parse(ps)
	n.n().To = ps.Pos()
	n.n().sourceText = ps.Src()[begin:ps.Pos()]
	return parsed[N]{n}
}

type parsed[N Node] struct {
	n N
}

func (p parsed[N]) addAs(ptr *N, parent Node) {
	*ptr = p.n
	addChild(parent, p.n)
}

func (p parsed[N]) addTo(ptr *[]N, parent Node) {
	*ptr = append(*ptr, p.n)
	addChild(parent, p.n)
}

func addChild(p Node, ch Node) {
	p.n().addChild(ch)
	ch.n().parent = p
}

// Tells the parser that parsing is done.
func (ps *parser) done() {
	if ps.Pos() != len(ps.Src()) {
		r, _ := utf8.DecodeRuneInString(ps.Src()[ps.Pos():])
		ps.error(fmt.Errorf("unexpected rune %q", r))
	}
}

const eof rune = EOF // Use the exported EOF constant from lexer

func (ps *parser) peek() rune {
	if ps.Pos() == len(ps.Src()) {
		return eof
	}
	r, _ := utf8.DecodeRuneInString(ps.Src()[ps.Pos():])
	return r
}

func (ps *parser) hasPrefix(prefix string) bool {
	return strings.HasPrefix(ps.Src()[ps.Pos():], prefix)
}

// next delegates to the embedded Lexer's Next method
func (ps *parser) next() rune {
	return ps.Next()
}

// backup delegates to the embedded Lexer's Backup method
func (ps *parser) backup() {
	ps.Backup()
}

func (ps *parser) errorp(r diag.Ranger, e error) {
	err := &Error{
		Message: e.Error(),
		Context: *diag.NewContext(ps.srcName, ps.Src(), r),
		Partial: r.Range().From == len(ps.Src()),
	}
	ps.errors = append(ps.errors, err)
}

func (ps *parser) error(e error) {
	end := ps.Pos()
	if end < len(ps.Src()) {
		end++
	}
	ps.errorp(diag.Ranging{From: ps.Pos(), To: end}, e)
}

// UnpackErrors returns the constituent parse errors if the given error contains
// one or more parse errors. Otherwise it returns nil.
func UnpackErrors(e error) []*Error {
	if errs := diag.UnpackErrors[ErrorTag](e); len(errs) > 0 {
		return errs
	}
	return nil
}

func newError(text string, shouldbe ...string) error {
	if len(shouldbe) == 0 {
		return errors.New(text)
	}
	var buf bytes.Buffer
	if len(text) > 0 {
		buf.WriteString(text + ", ")
	}
	buf.WriteString("should be " + shouldbe[0])
	for i, opt := range shouldbe[1:] {
		if i == len(shouldbe)-2 {
			buf.WriteString(" or ")
		} else {
			buf.WriteString(", ")
		}
		buf.WriteString(opt)
	}
	return errors.New(buf.String())
}
