package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/socialplanner/instahelper/app/update"
)

// APIUpdateHandler is the handler used to update instahelper to the latest version
func APIUpdateHandler(w http.ResponseWriter, r *http.Request) {

}

// APIUpdateToHandler is the handler used to update instahelper to a specific version
func APIUpdateToHandler(w http.ResponseWriter, r *http.Request) {
	ver := chi.URLParam(r, "version")

	// No version passed
	if ver == "" {
		http.Error(w, "Invalid Version", http.StatusBadRequest)
		return
	}

	asset, err := update.To(ver)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	b, _ := json.Marshal(asset)

	w.Header().Set("Content-Type", "application/json")
	w.Write(b)
}
