package templum

import (
	"fmt"
	"html/template"
	"io/fs"
	"maps"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing/fstest"

	"gopkg.in/yaml.v3"
)

type (
	MemoryFS   = fstest.MapFS
	MemoryFile = fstest.MapFile
)

type SiteContext struct {
	BaseURL    string
	Config     map[string]string
	Pages      []*Page
	OtherFiles fs.FS
}

type PageContext struct {
	*SiteContext

	Page    *Page
	Content template.HTML
}

func newSiteContext(baseURL, contentPath string) (*SiteContext, error) {
	fsys := os.DirFS(contentPath)

	config, err := config(fsys)
	if err != nil {
		return nil, fmt.Errorf("error reading config: %w", err)
	}

	pages, otherFiles, err := generatePages(fsys, ".")
	if err != nil {
		return nil, fmt.Errorf("error reading docs: %w", err)
	}

	return &SiteContext{
		BaseURL:    baseURL,
		Config:     config,
		Pages:      pages,
		OtherFiles: otherFiles,
	}, nil
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

func generatePages(fsys fs.FS, root string) ([]*Page, MemoryFS, error) {
	var pages []*Page

	entries, err := fs.ReadDir(fsys, root)
	if err != nil {
		return nil, nil, err
	}

	otherFiles := map[string]*MemoryFile{}

	for _, entry := range entries {
		if strings.HasPrefix(entry.Name(), "_") {
			continue
		}

		path := filepath.Join(root, entry.Name())

		switch {
		case entry.IsDir():
			section := NewSectionPage(path)

			subPages, subOtherFiles, err := generatePages(fsys, path)
			if err != nil {
				return nil, nil, err
			}

			maps.Copy(otherFiles, subOtherFiles)

			section.AddChildren(subPages...)

			pages = append(pages, section)
		case filepath.Ext(entry.Name()) == ".md":
			pages = append(pages, NewMarkdownPage(fsys, path))
		default:
			data, err := fs.ReadFile(fsys, path)
			if err != nil {
				return nil, nil, err
			}

			otherFiles[slug(path)+filepath.Ext(path)] = &MemoryFile{Data: data}
		}
	}

	slices.SortFunc(pages, sortPages)

	return pages, otherFiles, nil
}

func sortPages(a, b *Page) int {
	if a.Order() == b.Order() {
		return strings.Compare(a.Title(), b.Title())
	}

	return a.Order() - b.Order()
}
