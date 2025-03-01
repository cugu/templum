package plain

import (
	"context"
	"io"

	"github.com/a-h/templ"

	"github.com/cugu/templum"
)

func script(c *templum.PageContext) templ.Component {
	return templ.ComponentFunc(func(_ context.Context, w io.Writer) error {
		s := "<script>var base_url = '" + c.BaseURL + "';</script>"

		_, err := io.WriteString(w, s)

		return err
	})
}
