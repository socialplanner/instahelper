package handlers

import (
	tpl "html/template"
	"net/http"
	"sort"
)

// Page represents an individual page for instahelper.
type Page struct {

	// How the pages should be arranged on the sidebar
	ID int

	// Doubles as Title
	Name string

	// Link on the sidebar
	Link string

	// Material Design Icons
	// Choose from https://material.io/icons/
	// Replace all spaces with underscores
	Icon string

	// html/template.Template to execute
	Template *tpl.Template

	Handler http.HandlerFunc

	// To show it on sidebar or not
	Unlisted bool
}

// SortPages will return the list of pages sorted by ID.
//
// See: https://stackoverflow.com/questions/18695346/how-to-sort-a-mapstringint-by-its-values
func SortPages(m map[string]Page) []Page {
	// Temp struct
	type kv struct {
		Key string
		Val Page
	}

	var pages []kv
	for k, v := range m {
		pages = append(pages, kv{k, v})
	}

	sort.Slice(pages, func(i, j int) bool {
		return pages[i].Val.ID < pages[j].Val.ID
	})

	out := []Page{}

	for _, val := range pages {
		out = append(out, val.Val)
	}
	return out
}
