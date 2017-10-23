package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"github.com/sirupsen/logrus"
	"github.com/socialplanner/instahelper/app/config"
	"github.com/socialplanner/instahelper/app/update"
)

// UpdateHandler is the handler for /update
func UpdateHandler(w http.ResponseWriter, r *http.Request) {

	var releases []update.Release
	var err error

	if r := config.Get("releases"); r == nil {
		releases, err = update.ListReleases()
	} else {
		if rel, ok := r.([]update.Release); ok {
			releases = rel
		} else {
			releases, err = update.ListReleases()
		}
	}

	if err != nil {
		logrus.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		Error(w, err)
		return
	}

	for index, r := range releases {
		// Ignore error to default to the raw date
		t, err := time.Parse(time.RFC3339, r.PublishedAt)

		if err == nil {
			releases[index].PublishedAt = t.Format("Mon Jan 2, 3:04 PM ")
		}
	}

	err = Template("update").Execute(w, map[string]interface{}{
		"Releases": releases,
	})

	if err != nil {
		logrus.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		Error(w, err)
	}
}

// APIUpdateHandler is the handler used to update instahelper to the latest version
func APIUpdateHandler(w http.ResponseWriter, r *http.Request) {
	asset, err := update.ToLatest(update.VERSION)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	b, _ := json.Marshal(asset)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(b)
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
	w.WriteHeader(http.StatusOK)
	w.Write(b)
}
