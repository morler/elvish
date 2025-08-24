package tk

import (
	"strings"
	"testing"

	"src.elv.sh/pkg/ui"
)

// BenchmarkLabelRender tests label rendering performance with large content
func BenchmarkLabelRender(b *testing.B) {
	tests := []struct {
		name     string
		content  ui.Text
		width    int
		height   int
	}{
		{
			name:    "short-content",
			content: ui.T("Hello world"),
			width:   80,
			height:  5,
		},
		{
			name:    "long-single-line",
			content: ui.T(strings.Repeat("Lorem ipsum dolor sit amet ", 100)),
			width:   80,
			height:  5,
		},
		{
			name:    "multiline-content",
			content: ui.T(strings.Repeat("Line content\n", 100)),
			width:   80,
			height:  10,
		},
		{
			name:    "large-multiline-height-limited",
			content: ui.T(strings.Repeat("Large content line with many characters\n", 1000)),
			width:   80,
			height:  20,
		},
	}

	for _, tt := range tests {
		b.Run(tt.name, func(b *testing.B) {
			label := Label{Content: tt.content}
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				_ = label.Render(tt.width, tt.height)
			}
		})
	}
}

// BenchmarkListBoxRender tests listbox rendering performance with many items
func BenchmarkListBoxRender(b *testing.B) {
	tests := []struct {
		name      string
		nItems    int
		width     int
		height    int
		multiline bool
	}{
		{
			name:   "few-items-vertical",
			nItems: 10,
			width:  80,
			height: 20,
		},
		{
			name:   "many-items-vertical",
			nItems: 1000,
			width:  80,
			height: 20,
		},
		{
			name:   "many-items-small-height",
			nItems: 1000,
			width:  80,
			height: 5,
		},
		{
			name:      "multiline-items",
			nItems:    500,
			width:     80,
			height:    25,
			multiline: true,
		},
	}

	for _, tt := range tests {
		b.Run(tt.name, func(b *testing.B) {
			var items TestItems
			if tt.multiline {
				items = TestItems{NItems: tt.nItems, Prefix: "Multi\nline\nitem "}
			} else {
				items = TestItems{NItems: tt.nItems}
			}
			
			listBox := NewListBox(ListBoxSpec{
				State: ListBoxState{Items: items, Selected: 0},
			})
			
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				_ = listBox.Render(tt.width, tt.height)
			}
		})
	}
}

// BenchmarkListBoxHorizontalRender tests horizontal listbox rendering performance
func BenchmarkListBoxHorizontalRender(b *testing.B) {
	tests := []struct {
		name   string
		nItems int
		width  int
		height int
	}{
		{
			name:   "few-items-horizontal",
			nItems: 10,
			width:  120,
			height: 10,
		},
		{
			name:   "many-items-horizontal",
			nItems: 500,
			width:  120,
			height: 10,
		},
	}

	for _, tt := range tests {
		b.Run(tt.name, func(b *testing.B) {
			items := TestItems{NItems: tt.nItems}
			listBox := NewListBox(ListBoxSpec{
				Horizontal: true,
				State:      ListBoxState{Items: items, Selected: 0},
			})
			
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				_ = listBox.Render(tt.width, tt.height)
			}
		})
	}
}