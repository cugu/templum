package templum

import (
	"context"
	"fmt"
	"io/fs"
	"os"
	"testing/fstest"

	"gopkg.in/yaml.v3"
)

func Generate(ctx context.Context, baseURL, contentPath string, theme Theme, outputPath string) error {
	fsys := os.DirFS(contentPath)

	config, err := config(fsys)
	if err != nil {
		return fmt.Errorf("error reading config: %w", err)
	}

	docs, err := fs.Sub(fsys, "docs")
	if err != nil {
		return fmt.Errorf("error reading docs: %w", err)
	}

	pages, otherFiles, err := newPages(docs, ".")
	if err != nil {
		return fmt.Errorf("error reading docs: %w", err)
	}

	staticFS, err := staticFS(fsys)
	if err != nil {
		return fmt.Errorf("error reading static: %w", err)
	}

	content := Content{
		BaseURL: baseURL,
		Pages:   pages,
		Config:  config,
	}

	if err := generate(ctx, content, theme, outputPath, []fs.FS{staticFS, fstest.MapFS(otherFiles)}); err != nil {
		return fmt.Errorf("error generating: %w", err)
	}

	return nil
}

func generate(ctx context.Context, content Content, theme Theme, outputPath string, fsyss []fs.FS) error {
	docFS, err := theme.Render(ctx, content)
	if err != nil {
		return err
	}

	return writeToDisk(append([]fs.FS{docFS}, fsyss...), outputPath)
}

func config(fsys fs.FS) (map[string]string, error) {
	var config map[string]string

	f, err := fsys.Open("config.yaml")
	if err != nil {
		return nil, err
	}
	defer f.Close()

	if err := yaml.NewDecoder(f).Decode(&config); err != nil {
		return nil, err
	}

	return config, nil
}

func staticFS(fsys fs.FS) (fs.FS, error) {
	static, err := fs.Sub(fsys, "static")
	if err != nil {
		return nil, err
	}

	out := fstest.MapFS{}

	err = fs.WalkDir(static, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() {
			return nil
		}

		b, err := fs.ReadFile(static, path)
		if err != nil {
			return err
		}

		out[path] = &fstest.MapFile{Data: b}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return out, nil
}
