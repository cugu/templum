package blog

import (
	"bytes"
	"context"
	"io/fs"
	"maps"

	"github.com/cugu/templum"
	"github.com/cugu/templum/theme/blog/static"
)

var _ templum.Theme = Theme{}

type Theme struct{}

func (t Theme) Render(ctx context.Context, context *templum.SiteContext) (fs.FS, error) {
	memoryFS := templum.MemoryFS{}

	files, err := toFiles(ctx, context, context.Pages)
	if err != nil {
		return nil, err
	}

	maps.Copy(memoryFS, files)

	memoryFS["style.css"] = &templum.MemoryFile{Data: static.CSS}

	return memoryFS, nil
}

func toFiles(ctx context.Context, context *templum.SiteContext, pages []*templum.Page) (map[string]*templum.MemoryFile, error) {
	files := map[string]*templum.MemoryFile{}

	for _, p := range pages {
		if p.Type() == templum.Markdown {
			memoryFile, err := toFile(ctx, context, p)
			if err != nil {
				return nil, err
			}

			files[p.Link()] = memoryFile
		}

		if p.Type() == templum.Section {
			subFiles, err := toFiles(ctx, context, p.Children())
			if err != nil {
				return nil, err
			}

			maps.Copy(files, subFiles)
		}
	}

	return files, nil
}

func toFile(ctx context.Context, context *templum.SiteContext, p *templum.Page) (*templum.MemoryFile, error) {
	mainContent, err := p.HTML()
	if err != nil {
		return nil, err
	}

	mainComponent := html(&templum.PageContext{SiteContext: context, Page: p}, mainContent)

	var htmlBuffer bytes.Buffer
	if err := mainComponent.Render(ctx, &htmlBuffer); err != nil {
		return nil, err
	}

	return &templum.MemoryFile{Data: htmlBuffer.Bytes()}, nil
}
