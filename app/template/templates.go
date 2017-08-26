package template

import (
	"html/template"
	"io"

	"github.com/socialplanner/instahelper/app/assets"
)

var tmpls = map[string]Page{
	// Main Dashboard
	"dashboard": {
		Name:     "Dashboard",
		Link:     "/",
		Icon:     "dashboard",
		Template: newTemplate("base.html", "dashboard.html"),
	},

	"register": {
		Name:     "Add Account",
		Link:     "/register",
		Icon:     "person_add",
		Template: newTemplate("base.html", "register.html"),
	},
}

var funcs = template.FuncMap{
	"notifications": func() []Notification {
		// TODO. Replace with a method to actually get notifications
		return []Notification{
			{
				Text: "Test Notification",
				Link: "https://twitter.com/spieswithin",
			},
			{
				Text: "Test Notification 2",
				Link: "https://twitter.com/spieswithin",
			},
		}
	},
}

// Template will load the corresponding template with presets.
func Template(templateName string) *Page {
	if page, ok := tmpls[templateName]; ok {
		return &page
	}
	return nil
}

var a = assets.MustAsset

// Creates a template with the default funcs. Panics on error.
func newTemplate(files ...string) *template.Template {
	tmpl := template.New("*").Funcs(funcs)

	for _, file := range files {
		// assets.Asset defaults to '/' as a separator
		file = "templates/" + file

		content := string(a(file))
		tmpl = template.Must(tmpl.Parse(content))
	}
	return tmpl
}

// Execute is shorthand for Page.Template.Execute(w, Page, data)
func (p *Page) Execute(w io.Writer, data ...map[string]interface{}) error {

	templateData := map[string]interface{}{
		"Pages": tmpls,
		"Title": p.Name,
		"Icon":  p.Icon,
		"Link":  p.Link,
	}

	if len(data) > 0 {
		for key, val := range data[0] {
			templateData[key] = val
		}
	}

	return p.Template.Execute(
		w,
		templateData,
	)
}
