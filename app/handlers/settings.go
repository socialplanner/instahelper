package handlers

import (
	"net/http"

	"github.com/asdine/storm"

	"github.com/sirupsen/logrus"

	"github.com/go-chi/chi"
	"github.com/socialplanner/instahelper/app/config"
)

// AccSettingsHandler is the handler for the page to change account settings
// WIP
func AccSettingsHandler(w http.ResponseWriter, r *http.Request) {
	p := Template("account")

	username := chi.URLParam(r, "username")

	var acc config.Account

	if err := config.DB.One("Username", username, &acc); err != nil {
		if err == storm.ErrNotFound {
			http.NotFound(w, r)
			return
		}

		Error(w, err)
		return
	}

	err := p.Execute(w, map[string]interface{}{
		"Acc": acc,
	})

	if err != nil {
		Error(w, err)
	}
}

// SettingsHandler is the handler for instahelper settings
func SettingsHandler(w http.ResponseWriter, r *http.Request) {
	p := Template("settings")

	c, err := config.Config()

	if err != nil {
		Error(w, err)
		return
	}

	err = p.Execute(w, map[string]interface{}{
		"C": c,
	})

	if err != nil {
		Error(w, err)
	}
}

// APISettingsEditHandler is the api handler for editing instahelper settings
func APISettingsEditHandler(w http.ResponseWriter, r *http.Request) {
	username := r.PostFormValue("username")
	password := r.PostFormValue("password")
	analytics := r.PostFormValue("analytics") == "on"

	c, err := config.Config()

	if err != nil {
		logrus.Error(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	c.Username = username

	// Type none for no password
	if password != "none" {
		c.Password = password
	} else {
		c.Password = ""
	}
	c.Analytics = analytics

	if err := c.Update(); err != nil {
		logrus.Error(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write([]byte("Updated!"))
}
