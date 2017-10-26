package handlers

import (
	"fmt"
	"net/http"
)

// APIAccSettingsHandler is the handler to change the account settings
func APIAccSettingsHandler(w http.ResponseWriter, r *http.Request) {
	username := r.PostFormValue("username")
	fullName := r.PostFormValue("fullName")
	password := r.PostFormValue("password")
	proxy := r.PostFormValue("proxy")
	bio := r.PostFormValue("biography")

	fmt.Println(username, fullName, password, proxy, bio)
	w.Write([]byte("I'm useless so far!"))
}
