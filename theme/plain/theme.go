package plain

import (
	"bytes"
	"context"
	"encoding/json"
	"io/fs"
	"maps"

	"github.com/cugu/templum"
	"github.com/cugu/templum/theme/plain/static"
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
	memoryFS["main.js"] = &templum.MemoryFile{Data: static.JS}
	memoryFS["search.js"] = &templum.MemoryFile{Data: []byte(string(searchJS(context.Pages)) + string(static.SearchJS))}

	return memoryFS, nil
}

func toFiles(ctx context.Context, context *templum.SiteContext, pages []*templum.Page) (map[string]*templum.MemoryFile, error) {
	files := map[string]*templum.MemoryFile{}

	for _, p := range pages {
		if p.Type() == templum.MarkdownPage {
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
	mainContent, err := p.Markdown()
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

func searchJS(pages []*templum.Page) []byte {
	b, _ := json.Marshal(searchIndex(pages))

	return []byte("var index = " + string(b) + ";")
}

func searchIndex(pages []*templum.Page) []map[string]string {
	var data []map[string]string

	for _, p := range pages {
		if p.Type() == templum.MarkdownPage {
			md, err := p.Markdown()
			if err != nil {
				continue
			}

			data = append(data, map[string]string{
				"title": p.Title(),
				"href":  p.Link(),
				"body":  md,
			})
		}

		if p.Type() == templum.Section {
			data = append(data, searchIndex(p.Children())...)
		}
	}

	return data
}
