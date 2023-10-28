package tempel

import (
	"io"
	"io/fs"
	"os"
	"path/filepath"
)

func writeToDisk(fsyss []fs.FS, outputPath string) error {
	// remove all files and folder in public
	entries, err := os.ReadDir(outputPath)
	if err == nil {
		for _, entry := range entries {
			os.RemoveAll(filepath.Join(outputPath, entry.Name()))
		}
	}

	if err := os.MkdirAll(outputPath, 0o755); err != nil {
		return err
	}

	for _, fsys := range fsyss {
		if err := fs.WalkDir(fsys, ".", func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				return err
			}

			if d.IsDir() {
				return nil
			}

			f, err := fsys.Open(path)
			if err != nil {
				return err
			}
			defer f.Close()

			_ = os.MkdirAll(filepath.Dir(filepath.Join(outputPath, path)), 0o755)

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
	}

	return nil
}
