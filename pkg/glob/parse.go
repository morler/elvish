package glob

import (
	"bytes"

	"src.elv.sh/pkg/parse"
)

// Parse parses a pattern.
func Parse(s string) Pattern {
	segments := []Segment{}
	add := func(seg Segment) {
		segments = append(segments, seg)
	}
	p := parse.NewLexer(s)

rune:
	for {
		r := p.Next()
		switch r {
		case parse.EOF:
			break rune
		case '?':
			add(Wild{Question, false, nil})
		case '*':
			n := 1
			for p.Next() == '*' {
				n++
			}
			p.Backup()
			if n == 1 {
				add(Wild{Star, false, nil})
			} else {
				add(Wild{StarStar, false, nil})
			}
		case '/':
			for p.Next() == '/' {
			}
			p.Backup()
			add(Slash{})
		case '\\':
			fallthrough
		default:
			var literal bytes.Buffer

			// Handle the initial character
			if r == '\\' {
				// Process escaped character
				r = p.Next()
				if r == parse.EOF {
					literal.WriteRune('\\')
				} else {
					literal.WriteRune(r)
				}
				r = p.Next()
			}

			// Continue processing literal characters
			for {
				switch r {
				case '?', '*', '/', parse.EOF:
					goto endLiteral
				case '\\':
					// Backslash is always an escape character in literals
					r = p.Next()
					if r == parse.EOF {
						goto endLiteral
					}
					literal.WriteRune(r)
				default:
					literal.WriteRune(r)
				}
				r = p.Next()
			}
		endLiteral:
			p.Backup()
			add(Literal{literal.String()})
		}
	}
	return Pattern{segments, ""}
}

// Code duplication with parse/parser.go successfully eliminated.
// This parser now uses the shared parse.Lexer for common text reading operations,
// removing the previously duplicated parser struct, eof constant, next() and backup() methods.
