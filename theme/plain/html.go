package plain

import (
	"context"
	"io"

	"github.com/a-h/templ"
)

func raw(s string) templ.Component {
	return templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
		_, err := io.WriteString(w, s)

		return err
	})
}

func script(config map[string]string) templ.Component {
	return templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
		baseURL := config["base_url"]

		s := "<script>var base_url = '" + baseURL + "';</script>"

		_, err := io.WriteString(w, s)

		return err
	})
}

func style(light, dark string) templ.Component {
	return templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
		s := "<style>"
		s += "html.dark {"
		s += dark
		s += "}"
		s += "html:not(.dark) {"
		s += light
		s += "}"
		s += "</style>"

		_, err := io.WriteString(w, s)

		return err
	})
}
