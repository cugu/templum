package plain

import (
	"bytes"
	"context"
	"encoding/json"
	"io/fs"
	"maps"
	"testing/fstest"

	"github.com/cugu/templum"
	"github.com/cugu/templum/theme/plain/static"
)

var _ templum.Theme = Theme{}

type Theme struct{}

type pageContext struct {
	Title  string
	Active string
}

func (t Theme) Render(ctx context.Context, content templum.Content) (fs.FS, error) {
	memoryFS := fstest.MapFS{}

	files, err := toFiles(ctx, content, content.Pages)
	if err != nil {
		return nil, err
	}

	maps.Copy(memoryFS, files)

	memoryFS["style.css"] = &fstest.MapFile{Data: static.CSS}
	memoryFS["main.js"] = &fstest.MapFile{Data: static.JS}
	memoryFS["accordion.js"] = &fstest.MapFile{Data: static.AccordionJS}
	memoryFS["search.js"] = &fstest.MapFile{Data: []byte(string(searchJS(content)) + string(static.SearchJS))}
	memoryFS["prism.js"] = &fstest.MapFile{Data: static.PrismJS}
	memoryFS["prism.css"] = &fstest.MapFile{Data: static.PrismCSS}
	memoryFS["prism-include-languages.js"] = &fstest.MapFile{Data: static.PrismIncludeLanguages}

	return memoryFS, nil
}

func searchJS(content templum.Content) []byte {
	data := searchIndex(content.Pages)

	b, _ := json.Marshal(data)

	return []byte("var index = " + string(b) + ";")
}

func searchIndex(pages []*templum.Page) []map[string]string {
	var data []map[string]string

	for _, p := range pages {
		if p.Type() == templum.Markdown {
			md, err := p.Markdown()
			if err != nil {
				continue
			}

			data = append(data, map[string]string{
				"title": p.Title(),
				"href":  p.Href(),
				"body":  md,
			})
		}

		if p.Type() == templum.Section {
			data = append(data, searchIndex(p.Children())...)
		}
	}

	return data
}

func toFiles(ctx context.Context, content templum.Content, pages []*templum.Page) (map[string]*fstest.MapFile, error) {
	files := map[string]*fstest.MapFile{}

	for _, p := range pages {
		if p.Type() == templum.Markdown {
			memoryFile, err := toFile(ctx, content, p)
			if err != nil {
				return nil, err
			}

			files[p.Href()] = memoryFile
		}

		if p.Type() == templum.Section {
			subFiles, err := toFiles(ctx, content, p.Children())
			if err != nil {
				return nil, err
			}

			maps.Copy(files, subFiles)
		}
	}

	return files, nil
}

func toFile(ctx context.Context, content templum.Content, p *templum.Page) (*fstest.MapFile, error) {
	mainContent, err := p.HTML()
	if err != nil {
		return nil, err
	}

	context := &pageContext{
		Title:  p.Title(),
		Active: p.Slug(),
	}

	var htmlBuffer bytes.Buffer
	if err := html(context, content.Config, content.Pages, mainContent).Render(ctx, &htmlBuffer); err != nil {
		return nil, err
	}

	return &fstest.MapFile{Data: htmlBuffer.Bytes()}, nil
}
