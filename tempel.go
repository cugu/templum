package tempel

import (
	"context"
	"io/fs"
)

type Content struct {
	Config map[string]any
	Docs   fs.FS
}

type Theme interface {
	Render(ctx context.Context, content Content) (fs.FS, error)
}
