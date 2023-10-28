package templum

import (
	"testing"
)

func TestPage_Slug(t *testing.T) {
	type fields struct {
		path string
	}

	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{"page", fields{path: "page.md"}, "page"},
		{"index", fields{path: "index.md"}, "index"},
		{"nested", fields{path: "nested/page.md"}, "nested/page"},
		{"nested index", fields{path: "nested/index.md"}, "nested/index"},
		{"nested nested", fields{path: "nested/nested/page.md"}, "nested/nested/page"},
		{"nested nested index", fields{path: "nested/nested/index.md"}, "nested/nested/index"},

		{"camel", fields{path: "2 Quick start"}, "quick-start"},
		{"camel page", fields{path: "2 Quick start/1 Installation.md"}, "quick-start/installation"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Page{
				path: tt.fields.path,
			}
			if got := p.Slug(); got != tt.want {
				t.Errorf("Slug() = %v, want %v", got, tt.want)
			}
		})
	}
}
