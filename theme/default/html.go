package _default

import (
	"context"
	"github.com/a-h/templ"
	"io"
)

func HTML(s string) templ.Component {
	return templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
		_, err := io.WriteString(w, s)
		return err
	})
}
