package templum

import (
	"bytes"
	"context"
	"io"

	d2 "github.com/FurqanSoftware/goldmark-d2"
	"github.com/a-h/templ"
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
)

func md2html(md string) (string, error) {
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
			&d2.Extender{
				Sketch: true,
			},
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

func Markdown(md string) templ.Component {
	return templ.ComponentFunc(func(_ context.Context, w io.Writer) error {
		s, err := md2html(md)
		if err != nil {
			return err
		}

		_, err = io.WriteString(w, s)

		return err
	})
}
