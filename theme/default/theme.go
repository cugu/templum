package _default

import (
	"bytes"
	"context"
	_ "embed"
	"encoding/json"
	"io/fs"
	"maps"
	"tempel"
	"testing/fstest"
)

//go:embed style.css
var css []byte

//go:embed main.js
var js []byte

//go:embed prism.js
var prismJS []byte

//go:embed prism.css
var prismCSS []byte

//go:embed prism-include-languages.js
var prismIncludeLanguages []byte

var _ tempel.Theme = DefaultTheme{}

type DefaultTheme struct{}

type PageContext struct {
	Active string
}

func (t DefaultTheme) Render(ctx context.Context, content tempel.Content) (fs.FS, error) {
	memoryFS := fstest.MapFS{}

	files, err := toFiles(ctx, content, content.Pages)
	if err != nil {
		return nil, err
	}

	maps.Copy(memoryFS, files)

	memoryFS["style.css"] = &fstest.MapFile{Data: css}
	memoryFS["search.js"] = &fstest.MapFile{Data: searchJS(content)}
	memoryFS["main.js"] = &fstest.MapFile{Data: js}
	memoryFS["prism.js"] = &fstest.MapFile{Data: prismJS}
	memoryFS["prism.css"] = &fstest.MapFile{Data: prismCSS}
	memoryFS["prism-include-languages.js"] = &fstest.MapFile{Data: prismIncludeLanguages}

	return memoryFS, nil
}

func searchJS(content tempel.Content) []byte {
	data := searchIndex(content.Pages)

	b, _ := json.Marshal(data)

	return []byte("var index = " + string(b) + ";")
}

func searchIndex(pages []*tempel.Page) []map[string]string {
	var data []map[string]string
	for _, p := range pages {
		if p.Type() == tempel.Markdown {
			md, err := p.Markdown()
			if err != nil {
				continue
			}

			data = append(data, map[string]string{
				"title": p.Title(),
				"href":  p.Href(),
				"body":  string(md),
			})
		}

		if p.Type() == tempel.Section {
			data = append(data, searchIndex(p.SubPages)...)
		}
	}

	return data
}

func toFiles(ctx context.Context, content tempel.Content, pages []*tempel.Page) (map[string]*fstest.MapFile, error) {
	files := map[string]*fstest.MapFile{}

	for _, p := range pages {
		if p.Type() == tempel.Markdown {
			memoryFile, err := toFile(ctx, content, p)
			if err != nil {
				return nil, err
			}

			files[p.Href()] = memoryFile
		}

		if p.Type() == tempel.Section {
			subFiles, err := toFiles(ctx, content, p.SubPages)
			if err != nil {
				return nil, err
			}

			maps.Copy(files, subFiles)
		}
	}

	return files, nil
}

func toFile(ctx context.Context, content tempel.Content, p *tempel.Page) (*fstest.MapFile, error) {
	html, err := p.HTML()
	if err != nil {
		return nil, err
	}

	context := &PageContext{
		Active: p.Slug(),
	}

	var htmlBuffer bytes.Buffer
	if err := page(context, content.Config, content.Pages, html).Render(ctx, &htmlBuffer); err != nil {
		return nil, err
	}

	return &fstest.MapFile{Data: htmlBuffer.Bytes()}, nil
}
