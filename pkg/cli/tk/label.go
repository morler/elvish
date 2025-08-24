package tk

import (
	"src.elv.sh/pkg/cli/term"
	"src.elv.sh/pkg/ui"
)

// Label is a Renderer that writes out a text.
type Label struct {
	Content ui.Text
}

// Render shows the content. If the given box is too small, the text is cropped.
func (l Label) Render(width, height int) *term.Buffer {
	return l.renderOptimized(width, height)
}

// MaxHeight returns the maximum height the Label can take when rendering within
// a bound box.
func (l Label) MaxHeight(width, height int) int {
	return len(l.render(width).Lines)
}

// render creates a buffer with the full content (used by MaxHeight)
func (l Label) render(width int) *term.Buffer {
	return term.NewBufferBuilder(width).WriteStyled(l.Content).Buffer()
}

// renderOptimized creates a buffer with height-aware early termination
func (l Label) renderOptimized(width, height int) *term.Buffer {
	bb := term.NewBufferBuilder(width)
	
	for _, seg := range l.Content {
		bb.WriteStringSGR(seg.Text, seg.Style.SGR())
		
		// Check if we've exceeded the height limit
		cursor := bb.Cursor()
		if cursor.Line >= height {
			// Trim to exact height and return
			buf := bb.Buffer()
			buf.TrimToLines(0, height)
			return buf
		}
	}
	
	return bb.Buffer()
}

// Handle always returns false.
func (l Label) Handle(event term.Event) bool {
	return false
}
