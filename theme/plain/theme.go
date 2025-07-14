package plain

import (
	"bytes"
	"context"
	"encoding/json"
	"html"
	"html/template"
	"io/fs"
	"maps"
	"strings"

	_ "embed"

	"github.com/cugu/templum"
	"github.com/cugu/templum/theme/plain/static"
)

var _ templum.Theme = Theme{}

type Theme struct{}

//go:embed layout.tmpl
var layoutTemplate string

var pageTemplate = template.Must(template.New("layout.tmpl").Funcs(template.FuncMap{
	"logo":            func(c *templum.PageContext) template.HTML { return logo(c) },
	"nav":             func(c *templum.PageContext) template.HTML { return navHTML(c) },
	"github_link":     func(c *templum.PageContext) template.HTML { return githubLink(c) },
	"menu_toggle":     func() template.HTML { return menuToggle() },
	"magnifyingGlass": func() template.HTML { return template.HTML(magnifyingGlass()) },
	"x":               func() template.HTML { return template.HTML(xIcon()) },
}).Parse(layoutTemplate))

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
	if style, ok := context.Config["style"]; ok {
		memoryFS["custom.css"] = &templum.MemoryFile{Data: []byte(string(style))}
	}

	return memoryFS, nil
}

func toFiles(ctx context.Context, context *templum.SiteContext, pages []*templum.Page) (map[string]*templum.MemoryFile, error) {
	files := map[string]*templum.MemoryFile{}

	for _, p := range pages {
		if p.Type() == templum.MarkdownPage {
			memoryFile, err := toFile(context, p)
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

func toFile(context *templum.SiteContext, p *templum.Page) (*templum.MemoryFile, error) {
	md, err := p.Markdown()
	if err != nil {
		return nil, err
	}

	htmlContent, err := templum.Markdown(&templum.PageContext{SiteContext: context, Page: p}, md)
	if err != nil {
		return nil, err
	}

	data := &templum.PageContext{SiteContext: context, Page: p, Content: htmlContent}

	var buf bytes.Buffer
	if err := pageTemplate.Execute(&buf, data); err != nil {
		return nil, err
	}

	return &templum.MemoryFile{Data: buf.Bytes()}, nil
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

// helper template functions
func navHTML(c *templum.PageContext) template.HTML {
	var sb strings.Builder
	writeNav(&sb, c, c.Pages, 0)

	return template.HTML(sb.String())
}

func writeNav(sb *strings.Builder, c *templum.PageContext, pages []*templum.Page, depth int) {
	sb.WriteString(`<nav class="flex flex-col space-y-2`)
	if depth > 0 {
		sb.WriteString(` pt-2`)
	}
	sb.WriteString(`">`)
	for _, p := range pages {
		if p.Type() == templum.Section {
			writeSection(sb, c, p, depth)
		} else {
			writeLink(sb, c, p)
		}
	}
	sb.WriteString(`</nav>`)
}

func writeSection(sb *strings.Builder, c *templum.PageContext, p *templum.Page, depth int) {
	sb.WriteString(`<details`)
	if strings.HasPrefix(c.Page.Link(), p.Link()) {
		sb.WriteString(` open`)
	}
	sb.WriteString(`>`)
	sb.WriteString(`<summary class="flex justify-between toggle block p-2 rounded cursor-pointer w-full hover:bg-gray-200 dark:hover:bg-gray-700`)
	if strings.HasPrefix(c.Page.Link(), p.Link()) {
		sb.WriteString(` active font-bold`)
	}
	sb.WriteString(`">`)
	sb.WriteString(html.EscapeString(p.Title()))
	sb.WriteString(`<div class="chevron">`)
	sb.WriteString(chevron())
	sb.WriteString(`</div></summary>`)
	if len(p.Children()) > 0 {
		sb.WriteString(`<div class="pl-1">`)
		writeNav(sb, c, p.Children(), depth+1)
		sb.WriteString(`</div>`)
	}
	sb.WriteString(`</details>`)
}

func writeLink(sb *strings.Builder, c *templum.PageContext, p *templum.Page) {
	sb.WriteString(`<a href="`)
	sb.WriteString(html.EscapeString(c.BaseURL + p.Link()))
	sb.WriteString(`" class="block p-2 rounded w-full hover:bg-gray-200 dark:hover:bg-gray-700`)
	if c.Page.Link() == p.Link() {
		sb.WriteString(` font-bold bg-gray-200 dark:bg-gray-700`)
	} else {
		sb.WriteString(` font-normal`)
	}
	sb.WriteString(`">`)
	sb.WriteString(html.EscapeString(p.Title()))
	sb.WriteString(`</a>`)
}

func logo(c *templum.PageContext) template.HTML {
	var sb strings.Builder
	sb.WriteString(`<a href="` + html.EscapeString(c.BaseURL) + `">`)
	sb.WriteString(`<div class="flex flex-row flex-1 my-1.5">`)
	if logoPath, ok := c.Config["logo"]; ok {
		sb.WriteString(`<img src="` + html.EscapeString(c.BaseURL+logoPath) + `" alt="logo" class="h-8 mr-2"/>`)
	}
	if title, ok := c.Config["title"]; ok {
		sb.WriteString(`<h1 class="text-2xl font-bold">` + html.EscapeString(title) + `</h1>`)
	}
	sb.WriteString(`</div></a>`)

	return template.HTML(sb.String())
}

func githubLink(c *templum.PageContext) template.HTML {
	url, ok := c.Config["github_url"]
	if !ok {
		return ""
	}
	var sb strings.Builder
	sb.WriteString(`<a href="` + html.EscapeString(url) + `" class="hover:bg-gray-100 dark:hover:bg-gray-700 border dark:border-gray-600 font-bold py-2 px-4 rounded-lg flex ml-2 flex-row space-x-2">`)
	sb.WriteString(github())
	sb.WriteString(`<span>GitHub</span></a>`)
	return template.HTML(sb.String())
}

func menuToggle() template.HTML {
	return template.HTML(`<div class="menu-toggle md:hidden cursor-pointer">` + hamburger() + `</div>`)
}

func github() string {
	return `<svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="feather feather-github"><path d="M9 19c-5 1.5-5-2.5-7-3m14 6v-3.87a3.37 3.37 0 0 0-.94-2.61c3.14-.35 6.44-1.54 6.44-7A5.44 5.44 0 0 0 20 4.77 5.07 5.07 0 0 0 19.91 1S18.73.65 16 2.48a13.38 13.38 0 0 0-7 0C6.27.65 5.09 1 5.09 1A5.07 5.07 0 0 0 5 4.77a5.44 5.44 0 0 0-1.5 3.78c0 5.42 3.3 6.61 6.44 7A3.37 3.37 0 0 0 9 18.13V22"></path></svg>`
}

func hamburger() string {
	return `<svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="feather feather-menu"><line x1="3" y1="12" x2="21" y2="12"></line><line x1="3" y1="6" x2="21" y2="6"></line><line x1="3" y1="18" x2="21" y2="18"></line></svg>`
}

func chevron() string {
	return `<svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="feather feather-chevron-right"><polyline points="9 18 15 12 9 6"></polyline></svg>`
}

func magnifyingGlass() string {
	return `<svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="feather feather-search"><circle cx="11" cy="11" r="8"></circle><line x1="21" y1="21" x2="16.65" y2="16.65"></line></svg>`
}

func xIcon() string {
	return `<svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="feather feather-x"><line x1="18" y1="6" x2="6" y2="18"></line><line x1="6" y1="6" x2="18" y2="18"></line></svg>`
}
