package handlers

import (
	"net/http"
	"strconv"

	"github.com/sirupsen/logrus"

	"github.com/go-chi/chi"
	"github.com/socialplanner/instahelper/app/config"
)

// AccSettingsHandler is the handler for the page to change account settings
// TODO
func AccSettingsHandler(w http.ResponseWriter, r *http.Request) {
	p := Template("account")

	username := chi.URLParam(r, "username")

	err := p.Execute(w, map[string]interface{}{
		"Username": username,
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
	portStr := r.PostFormValue("port")
	analytics := r.PostFormValue("analytics") == "on"

	c, err := config.Config()

	if err != nil {
		logrus.Error(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if portStr != "" {
		port, err := strconv.Atoi(portStr)

		if err != nil {
			logrus.Error(err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		c.Port = port
	}

	// If arg isn't the zero value. Replace the config with it.
	switch {
	case username != "":
		c.Username = username
		fallthrough
	case password != "":
		c.Password = password
	}

	c.Analytics = analytics

	if err := c.Update(); err != nil {
		logrus.Error(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write([]byte("Updated!"))
}
