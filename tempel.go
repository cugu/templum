package tempel

import (
	"context"
	"io/fs"
)

type Content struct {
	Config map[string]string
	Static fs.FS
	Pages  []*Page
}

type Theme interface {
	Render(ctx context.Context, content Content) (fs.FS, error)
}

type PageType int

const (
	Section PageType = iota
	Markdown
)
