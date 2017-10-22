package handlers

import (
	"net/http"

	"github.com/go-chi/chi"
)

// SettingsHandler is the handler for the page to change account settings
// TODO
func SettingsHandler(w http.ResponseWriter, r *http.Request) {
	p := Template("settings")

	username := chi.URLParam(r, "username")

	err := p.Execute(w, map[string]interface{}{
		"Username": username,
	})

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
