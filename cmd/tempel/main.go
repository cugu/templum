package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"tempel"
	_default "tempel/theme/default"
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
	docs, err := fs.Sub(fsys, "docs")
	if err != nil {
		return err
	}

	content := tempel.Content{
		Config: map[string]interface{}{},
		Docs:   docs,
	}

	t := _default.DefaultTheme{}
	out, err := t.Render(ctx, content)
	if err != nil {
		return err
	}

	return writeToDisk(out, outputPath)
}

func writeToDisk(out fs.FS, outputPath string) error {
	// remove all files and folder in public
	entries, err := os.ReadDir(outputPath)
	if err == nil {
		for _, entry := range entries {
			os.RemoveAll(filepath.Join(outputPath, entry.Name()))
		}
	}

	if err := os.MkdirAll(outputPath, 0755); err != nil {
		return err
	}

	if err := fs.WalkDir(out, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() {
			return nil
		}

		f, err := out.Open(path)
		if err != nil {
			return err
		}
		defer f.Close()

		os.MkdirAll(filepath.Dir(filepath.Join(outputPath, path)), 0755)

		out, err := os.Create(filepath.Join(outputPath, path))
		if err != nil {
			return err
		}
		defer out.Close()
		_, err = io.Copy(out, f)
		return err
	}); err != nil {
		return err
	}

	return nil
}
