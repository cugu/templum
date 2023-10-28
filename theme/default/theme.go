package _default

import (
	"bytes"
	"context"
	_ "embed"
	"fmt"
	"io/fs"
	"path/filepath"
	"slices"
	"strings"
	"tempel"
	"testing/fstest"

	"github.com/yuin/goldmark"
)

var _ tempel.Theme = DefaultTheme{}

type DefaultTheme struct{}

//go:embed out.css
var css []byte

//go:embed out.js
var js []byte

type Link struct {
	ID   string
	Href string
	Text string

	Children []Link
}

func tempelNav(content tempel.Content, root string) ([]Link, error) {
	var nav []Link

	entries, err := fs.ReadDir(content.Docs, root)
	if err != nil {
		return nil, err
	}

	for _, entry := range entries {
		path := filepath.Join(root, entry.Name())

		link := Link{
			ID:   path,
			Text: navText(path),
		}

		if entry.IsDir() {
			children, err := tempelNav(content, path)
			if err != nil {
				return nil, err
			}

			link.Children = children
		} else {
			link.Href = path[:len(path)-len(filepath.Ext(path))] + ".html"
		}

		nav = append(nav, link)
	}

	slices.SortFunc(nav, func(a, b Link) int {
		// index.md should always be first
		if a.ID == "index.md" {
			return -1
		}

		return strings.Compare(a.ID, b.ID)
	})

	return nav, nil
}

func navText(path string) string {
	base := filepath.Base(path)

	if base == "index.md" {
		return "Home"
	}

	// remove prefix
	for i := 0; i < len(base); i++ {
		if base[i] == '-' {
			base = base[i+1:]
			break
		}
	}

	// remove extension
	base = base[:len(base)-len(filepath.Ext(base))]

	return base
}

func (t DefaultTheme) Render(ctx context.Context, content tempel.Content) (fs.FS, error) {
	fmt.Println("Render called")

	memoryFS := fstest.MapFS{}

	nav, err := tempelNav(content, ".")
	if err != nil {
		return nil, err
	}

	if err := fs.WalkDir(content.Docs, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() || filepath.Ext(path) != ".md" {
			return nil
		}

		b, err := fs.ReadFile(content.Docs, path)
		if err != nil {
			return err
		}

		var contentBuffer, htmlBuffer bytes.Buffer

		if err := goldmark.Convert(b, &contentBuffer); err != nil {
			return err
		}

		if err := page(content.Config, nav, contentBuffer.String()).Render(ctx, &htmlBuffer); err != nil {
			return err
		}

		htmlPath := path[:len(path)-len(filepath.Ext(path))] + ".html"
		memoryFS[htmlPath] = &fstest.MapFile{Data: htmlBuffer.Bytes()}

		return nil
	}); err != nil {
		return nil, err
	}

	memoryFS["out.css"] = &fstest.MapFile{Data: css}
	memoryFS["out.js"] = &fstest.MapFile{Data: js}

	return memoryFS, nil
}
