package handlers

import (
	"fmt"
	"net/http"

	"github.com/asdine/storm"
	"github.com/go-chi/chi"
	"github.com/socialplanner/instahelper/app/config"
)

// APIAccProfileHandler is the handler to change the account profile settings
func APIAccProfileHandler(w http.ResponseWriter, r *http.Request) {
	username := r.PostFormValue("username")
	fullName := r.PostFormValue("fullName")
	password := r.PostFormValue("password")
	proxy := r.PostFormValue("proxy")
	bio := r.PostFormValue("biography")

	fmt.Println(username, fullName, password, proxy, bio)
	w.Write([]byte("I'm useless so far!"))
}

// AccSettingsHandler is the handler for the page to change account settings
// To add a category, edit the page and add the template html.
// Make sure the form has the class jobs-settings
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
