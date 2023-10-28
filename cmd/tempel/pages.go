package main

import (
	"io/fs"
	"path/filepath"
	"slices"
	"strings"

	"github.com/cugu/tempel"
)

func newPages(fsys fs.FS, root string) ([]*tempel.Page, error) {
	var pages []*tempel.Page

	entries, err := fs.ReadDir(fsys, root)
	if err != nil {
		return nil, err
	}

	for _, entry := range entries {
		if !(entry.IsDir() || filepath.Ext(entry.Name()) == ".md") {
			continue
		}

		path := filepath.Join(root, entry.Name())

		switch {
		case entry.IsDir():
			section := tempel.NewSectionPage(path)

			subPages, err := newPages(fsys, path)
			if err != nil {
				return nil, err
			}

			section.SubPages = subPages

			pages = append(pages, section)
		case filepath.Ext(entry.Name()) == ".md":
			pages = append(pages, tempel.NewMarkdownPage(fsys, path))
		default:
			continue
		}
	}

	slices.SortFunc(pages, sortPages)

	return pages, nil
}

func sortPages(a, b *tempel.Page) int {
	if a.Order() == b.Order() {
		return strings.Compare(a.Title(), b.Title())
	}

	return a.Order() - b.Order()
}
