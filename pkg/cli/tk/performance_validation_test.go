package tk

import (
	"strings"
	"testing"

	"src.elv.sh/pkg/ui"
)

// Test to validate that our optimizations produce identical results to the original approach
func TestLabelOptimizedRenderEquivalence(t *testing.T) {
	testCases := []struct {
		name    string
		content ui.Text
		width   int
		height  int
	}{
		{
			name:    "short-content",
			content: ui.T("Hello world"),
			width:   80,
			height:  5,
		},
		{
			name:    "multiline-content",
			content: ui.T(strings.Repeat("Line content\n", 20)),
			width:   40,
			height:  10,
		},
		{
			name:    "long-content-cropped",
			content: ui.T(strings.Repeat("Very long line that will wrap multiple times\n", 50)),
			width:   20,
			height:  5,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			label := Label{Content: tc.content}

			// Get optimized result
			optimizedResult := label.renderOptimized(tc.width, tc.height)

			// Get unoptimized result (simulating the old approach)
			unoptimizedFull := label.render(tc.width)
			unoptimizedFull.TrimToLines(0, tc.height)

			// Compare results
			if optimizedResult.Width != unoptimizedFull.Width {
				t.Errorf("Width mismatch: optimized=%d, unoptimized=%d",
					optimizedResult.Width, unoptimizedFull.Width)
			}

			if len(optimizedResult.Lines) != len(unoptimizedFull.Lines) {
				t.Errorf("Line count mismatch: optimized=%d, unoptimized=%d",
					len(optimizedResult.Lines), len(unoptimizedFull.Lines))
			}

			// Compare line by line (up to the height limit)
			minLines := len(optimizedResult.Lines)
			if len(unoptimizedFull.Lines) < minLines {
				minLines = len(unoptimizedFull.Lines)
			}

			for i := 0; i < minLines; i++ {
				if len(optimizedResult.Lines[i]) != len(unoptimizedFull.Lines[i]) {
					t.Errorf("Line %d length mismatch: optimized=%d, unoptimized=%d",
						i, len(optimizedResult.Lines[i]), len(unoptimizedFull.Lines[i]))
				}

				for j := 0; j < len(optimizedResult.Lines[i]) && j < len(unoptimizedFull.Lines[i]); j++ {
					if optimizedResult.Lines[i][j] != unoptimizedFull.Lines[i][j] {
						t.Errorf("Cell mismatch at line %d, col %d", i, j)
					}
				}
			}
		})
	}
}

// Test that validates listbox optimization produces correct results
func TestListBoxOptimizedRenderCorrectness(t *testing.T) {
	testCases := []struct {
		name      string
		nItems    int
		width     int
		height    int
		multiline bool
	}{
		{
			name:   "few-items-fits",
			nItems: 5,
			width:  80,
			height: 10,
		},
		{
			name:   "many-items-height-limited",
			nItems: 100,
			width:  80,
			height: 5,
		},
		{
			name:      "multiline-items-cropped",
			nItems:    20,
			width:     40,
			height:    8,
			multiline: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var items TestItems
			if tc.multiline {
				items = TestItems{NItems: tc.nItems, Prefix: "Multi\nline\nitem "}
			} else {
				items = TestItems{NItems: tc.nItems}
			}

			listBox := NewListBox(ListBoxSpec{
				State: ListBoxState{Items: items, Selected: 0},
			})

			// Render with optimization
			result := listBox.Render(tc.width, tc.height)

			// Basic validation - should not exceed height
			if len(result.Lines) > tc.height {
				t.Errorf("Result exceeds height limit: got %d lines, max %d",
					len(result.Lines), tc.height)
			}

			// Width should match
			if result.Width != tc.width {
				t.Errorf("Width mismatch: got %d, expected %d", result.Width, tc.width)
			}

			// Should have some content unless no items
			if tc.nItems > 0 && len(result.Lines) == 0 {
				t.Error("No content rendered despite having items")
			}
		})
	}
}
