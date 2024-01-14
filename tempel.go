package templum

import (
	"context"
	"io/fs"
)

type Content struct {
	BaseURL string
	Config  map[string]string
	Pages   []*Page
}

type Theme interface {
	Render(ctx context.Context, content Content) (fs.FS, error)
}

type PageType int

const (
	Section PageType = iota
	Markdown
)
