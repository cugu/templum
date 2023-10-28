package plain

import (
	"context"
	"io"

	"github.com/a-h/templ"
)

func HTML(s string) templ.Component {
	return templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
		_, err := io.WriteString(w, s)

		return err
	})
}
