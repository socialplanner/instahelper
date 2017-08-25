package main

import (
	"fmt"
	"mime"
	"net/http"
	"path/filepath"

	"github.com/socialplanner/instahelper/app/assets"
	tmpl "github.com/socialplanner/instahelper/app/template"
)

func dashboard(w http.ResponseWriter, r *http.Request) {
	var err error
	url := r.URL.String()[1:]

	if url == "" {
		err = tmpl.Template("dashboard").Execute(w)
	} else {
		a, err := assets.Asset(url)
		if err == nil {
			w.Header().Set("Content-Type", mime.TypeByExtension(filepath.Ext(url)))
			fmt.Fprint(w, string(a))
		}
	}

	if err != nil {
		fmt.Fprint(w, err)
	}
}

func main() {
	fmt.Println("Rome wasn't built in a day.")
	http.HandleFunc("/", dashboard)
	http.ListenAndServe(":8080", nil)
}
