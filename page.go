package templum

import (
	"errors"
	"io/fs"
	"path"
	"path/filepath"
	"strconv"
	"strings"
)

type PageType int

const (
	Section PageType = iota
	MarkdownPage
)

type Page struct {
	fsys     fs.FS
	path     string
	pageType PageType

	children []*Page
}

func NewMarkdownPage(fsys fs.FS, path string) *Page {
	return &Page{
		fsys:     fsys,
		path:     path,
		pageType: MarkdownPage,
	}
}

func NewSectionPage(path string) *Page {
	return &Page{
		path:     path,
		pageType: Section,
	}
}

func (p *Page) Type() PageType {
	return p.pageType
}

func (p *Page) Title() string {
	base := filepath.Base(p.path)
	if base == "index.md" {
		return "Home"
	}

	filename, _ := baseParts(p.path)

	return filename
}

func (p *Page) Order() int {
	base := filepath.Base(p.path)
	if base == "index.md" {
		return 0
	}

	_, o := baseParts(p.path)

	return o
}

func baseParts(path string) (string, int) {
	base := filepath.Base(path)
	filename := base[:len(base)-len(filepath.Ext(base))]
	prefix, suffix, hasSpace := strings.Cut(filename, " ")

	if hasSpace {
		o, err := strconv.Atoi(prefix)
		if err != nil {
			return suffix, -1
		}

		return suffix, o
	}

	return filename, -1
}

func (p *Page) Link() string {
	if p.pageType == Section {
		return slug(p.path) + "/"
	}

	if p.path == "index.md" {
		return "index.html"
	}

	return slug(strings.TrimSuffix(p.path, filepath.Ext(p.path))) + ".html"
}

func (p *Page) Markdown() (string, error) {
	if p.pageType != MarkdownPage {
		return "", errors.New("page is not markdown")
	}

	b, err := fs.ReadFile(p.fsys, p.path)
	if err != nil {
		return "", err
	}

	return string(b), nil
}

func (p *Page) AddChildren(child ...*Page) {
	p.children = append(p.children, child...)
}

func (p *Page) Children() []*Page {
	return p.children
}

func slug(p string) string {
	htmlPath := ""

	for _, r := range strings.Split(p, "/") {
		name, _ := baseParts(r)
		htmlPath = path.Join(htmlPath, name)
	}

	for _, r := range []string{"\\", " ", ".", "_"} {
		htmlPath = strings.ReplaceAll(htmlPath, r, "-")
	}

	htmlPath = strings.ToLower(htmlPath)
	htmlPath = strings.Trim(htmlPath, "-")

	return htmlPath
}
