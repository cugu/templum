package tempel

import (
	"context"
	"io/fs"
)

type Content struct {
	Config map[string]any
	Pages  []*Page
}

type Theme interface {
	Render(ctx context.Context, content Content) (fs.FS, error)
}

type pageType int

const (
	Section pageType = iota
	Markdown
)
