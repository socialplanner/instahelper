package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"github.com/sirupsen/logrus"
	"github.com/socialplanner/instahelper/app/update"
)

// UpdateHandler is the handler for /update
func UpdateHandler(w http.ResponseWriter, r *http.Request) {
	releases, err := update.ListReleases()

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
			releases[index].PublishedAt = t.Format("Jan 2")
		}
	}

	err = Template("update").Execute(w, map[string]interface{}{
		"Releases":         releases,
		"Version":          update.VERSION,
		"DifferentVersion": update.DIFFERENTVERSION,
		"updateAvailable": func() bool {
			if asset, err := update.HigherVersion(update.VERSION); err == nil && asset != nil {
				return true
			}

			return false
		},

		"truncate": func(str string) string {
			truncateAt := 140
			var truncated string

			if len(str) > truncateAt {
				truncated = str[0:truncateAt]
			} else {
				return str
			}

			if truncated[len(truncated)-1:] != "." {
				truncated += "..."
			} else {
				truncated += ".."
			}

			return truncated
		},
	})

	if err != nil {
		logrus.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		Error(w, err)
	}
}

// APIUpdateHandler is the handler used to update instahelper to the latest version
func APIUpdateHandler(w http.ResponseWriter, r *http.Request) {
	_, err := update.ToLatest(update.VERSION)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "Awesome! Updated to the latest version! All you need to do is restart the current running app.")
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
