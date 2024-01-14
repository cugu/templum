package templum

import (
	"context"
	"fmt"
	"io/fs"
)

func Generate(ctx context.Context, baseURL, contentPath string, theme Theme, outputPath string) error {
	siteContext, err := newSiteContext(baseURL, contentPath)
	if err != nil {
		return fmt.Errorf("error creating site context: %w", err)
	}

	docFS, err := theme.Render(ctx, siteContext)
	if err != nil {
		return err
	}

	return writeToDisk(append([]fs.FS{docFS}, siteContext.OtherFiles), outputPath)
}
