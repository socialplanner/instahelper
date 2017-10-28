package handlers

import (
	"net/http"

	"github.com/sirupsen/logrus"

	"github.com/socialplanner/instahelper/app/config"
)

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
	updates := r.PostFormValue("automaticUpdates") == "on"

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
	c.AutomaticUpdates = updates

	if err := c.Update(); err != nil {
		logrus.Error(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write([]byte("Updated!"))
}
