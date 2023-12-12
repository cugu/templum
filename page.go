package templum

import (
	"bytes"
	"errors"
	"io/fs"
	"path"
	"path/filepath"
	"strconv"
	"strings"

	chromahtml "github.com/alecthomas/chroma/v2/formatters/html"
	"github.com/alecthomas/chroma/v2/styles"
	"github.com/stefanfritsch/goldmark-fences"
	"github.com/yuin/goldmark"
	highlighting "github.com/yuin/goldmark-highlighting/v2"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer/html"
	"go.abhg.dev/goldmark/anchor"
	"go.abhg.dev/goldmark/mermaid"
	"mvdan.cc/xurls/v2"
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
		pageType: Markdown,
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

func (p *Page) Slug() string {
	if p.path == "index.md" {
		return "index"
	}

	noExt := strings.TrimSuffix(p.path, filepath.Ext(p.path))

	htmlPath := ""

	for _, r := range strings.Split(noExt, "/") {
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

func (p *Page) Href() string {
	return p.Slug() + ".html"
}

func (p *Page) Markdown() (string, error) {
	if p.pageType != Markdown {
		return "", errors.New("page is not markdown")
	}

	b, err := fs.ReadFile(p.fsys, p.path)
	if err != nil {
		return "", err
	}

	return string(b), nil
}

func (p *Page) HTML() (string, string, string, error) {
	b, err := p.Markdown()
	if err != nil {
		return "", "", "", err
	}

	var htmlBuffer bytes.Buffer

	markdown := goldmark.New(
		goldmark.WithParserOptions(
			parser.WithAutoHeadingID(),
		),
		goldmark.WithRendererOptions(
			html.WithXHTML(),
			html.WithUnsafe(),
		),
		goldmark.WithExtensions(
			&anchor.Extender{
				Texter: anchor.Text("#"),
			},
			extension.NewLinkify(
				extension.WithLinkifyAllowedProtocols([][]byte{
					[]byte("http:"),
					[]byte("https:"),
				}),
				extension.WithLinkifyURLRegexp(
					xurls.Strict(),
				),
			),
			&mermaid.Extender{},
			&fences.Extender{},
			highlighting.NewHighlighting(
				highlighting.WithStyle("vs"),
				highlighting.WithFormatOptions(
					chromahtml.WithLineNumbers(true),
					chromahtml.WithClasses(true),
				),
			),
		),
	)

	if err := markdown.Convert([]byte(b), &htmlBuffer); err != nil {
		return "", "", "", err
	}

	formatter := chromahtml.New()

	light := &bytes.Buffer{}
	if err := formatter.WriteCSS(light, styles.Get("vs")); err != nil {
		return "", "", "", err
	}

	dark := &bytes.Buffer{}
	if err := formatter.WriteCSS(dark, styles.Get("github-dark")); err != nil {
		return "", "", "", err
	}

	return htmlBuffer.String(), light.String(), dark.String(), nil
}

func (p *Page) AddChildren(child ...*Page) {
	p.children = append(p.children, child...)
}

func (p *Page) Children() []*Page {
	return p.children
}
