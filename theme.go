package templum

import (
	"context"
	"io/fs"
)

type Theme interface {
	Render(ctx context.Context, siteContext *SiteContext) (fs.FS, error)
}
