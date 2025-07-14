package templum

import (
	"bytes"
	"html/template"

	d2 "github.com/FurqanSoftware/goldmark-d2"
	chromahtml "github.com/alecthomas/chroma/v2/formatters/html"
	fences "github.com/stefanfritsch/goldmark-fences"
	"github.com/yuin/goldmark"
	highlighting "github.com/yuin/goldmark-highlighting/v2"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer/html"
	"go.abhg.dev/goldmark/anchor"
	"go.abhg.dev/goldmark/mermaid"
	"mvdan.cc/xurls/v2"
	"oss.terrastruct.com/d2/d2layouts/d2elklayout"
	"oss.terrastruct.com/d2/d2themes/d2themescatalog"
)

func md2html(config map[string]string, md string) (string, error) {
	var htmlBuffer bytes.Buffer

	d2Settings := &d2.Extender{}

	if d2Sketch, ok := config["d2_sketch"]; ok {
		d2Settings.Sketch = d2Sketch == "true"
	}

	if d2ThemeName, ok := config["d2_theme_name"]; ok {
		for _, theme := range d2themescatalog.LightCatalog {
			if theme.Name == d2ThemeName {
				d2Settings.ThemeID = &theme.ID

				break
			}
		}

		for _, theme := range d2themescatalog.DarkCatalog {
			if theme.Name == d2ThemeName {
				d2Settings.ThemeID = &theme.ID

				break
			}
		}
	}

	if d2Layout, ok := config["d2_layout"]; ok {
		if d2Layout == "elk" {
			d2Settings.Layout = d2elklayout.DefaultLayout
		}
	}

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
			d2Settings,
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

	if err := markdown.Convert([]byte(md), &htmlBuffer); err != nil {
		return "", err
	}

	return htmlBuffer.String(), nil
}

func Markdown(c *PageContext, md string) (template.HTML, error) {
	s, err := md2html(c.Config, md)
	if err != nil {
		return "", err
	}

	return template.HTML(s), nil
}
