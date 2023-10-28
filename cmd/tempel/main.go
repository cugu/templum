package main

import (
	"context"
	"flag"
	"fmt"
	"gopkg.in/yaml.v3"
	"io/fs"
	"os"
	"tempel"
	_default "tempel/theme/default"
	"testing/fstest"
)

func main() {
	generateCmd(context.Background(), os.Args[1:])
}

func generateCmd(ctx context.Context, args []string) {
	cmd := flag.NewFlagSet("generate", flag.ExitOnError)
	contentPathFlag := cmd.String("content", ".", "filesystem path to content directory (default: '.')")
	outputPathFlag := cmd.String("output", "public", "filesystem path to output directory (default: 'public')")
	helpFlag := cmd.Bool("help", false, "Print help and exit.")
	err := cmd.Parse(args)
	if err != nil || *helpFlag {
		cmd.PrintDefaults()
		return
	}

	err = generate(ctx, *contentPathFlag, *outputPathFlag)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
}

func generate(ctx context.Context, contentPath, outputPath string) error {
	fsys := os.DirFS(contentPath)

	config, err := config(fsys)
	if err != nil {
		return err
	}

	docFS, err := docFS(ctx, fsys, config)
	if err != nil {
		return err
	}

	staticFS, err := staticFS(fsys)
	if err != nil {
		return err
	}

	return writeToDisk([]fs.FS{docFS, staticFS}, outputPath)
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

func docFS(ctx context.Context, fsys fs.FS, config map[string]string) (fs.FS, error) {
	docs, err := fs.Sub(fsys, "docs")
	if err != nil {
		return nil, err
	}

	pages, err := newPages(docs, ".")
	if err != nil {
		return nil, err
	}

	content := tempel.Content{
		Pages:  pages,
		Config: config,
	}

	t := _default.DefaultTheme{}
	docFS, err := t.Render(ctx, content)
	if err != nil {
		return nil, err
	}
	return docFS, nil
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
