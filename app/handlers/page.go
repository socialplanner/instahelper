package handlers

import tpl "html/template"

// Page represents an individual page for instahelper.
type Page struct {
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
}
