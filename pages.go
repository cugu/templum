package templum

import (
	"io/fs"
	"maps"
	"path/filepath"
	"slices"
	"strings"
	"testing/fstest"
)

func newPages(fsys fs.FS, root string) ([]*Page, map[string]*fstest.MapFile, error) {
	var pages []*Page

	entries, err := fs.ReadDir(fsys, root)
	if err != nil {
		return nil, nil, err
	}

	otherFiles := map[string]*fstest.MapFile{}

	for _, entry := range entries {
		path := filepath.Join(root, entry.Name())

		switch {
		case entry.IsDir():
			section := NewSectionPage(path)

			subPages, subOtherFiles, err := newPages(fsys, path)
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

			otherFiles[slug(path)+filepath.Ext(path)] = &fstest.MapFile{Data: data}
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
