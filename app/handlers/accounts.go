package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"github.com/socialplanner/instahelper/app/config"
	"github.com/socialplanner/instahelper/app/insta"
)

// APICreateAccountHandler is a http.Handler which should be used to save an account to db.
func APICreateAccountHandler(w http.ResponseWriter, r *http.Request) {
	// TODO, check if is a valid account/needs to break captcha.
	// TODO encrypt
	username := r.PostFormValue("username")
	password := r.PostFormValue("password")
	proxy := r.PostFormValue("proxy")

	if username == "" || password == "" {
		http.Error(w, "Invalid Form Input", http.StatusBadRequest)
		return
	}
	acc := &config.Account{}

	config.DB.One("Username", username, acc)

	if acc.Username == username {
		http.Error(w, "An account with that username already exists.", http.StatusConflict)
		return
	}

	// Encrypt password
	passwordENC, err := insta.Encrypt(password)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	ig, err := insta.CachedInsta(username, password, proxy)

	if err != nil {
		// Bad username/password combo or captcha
		http.Error(w, err.Error(), http.StatusBadGateway)
		return
	}

	err = config.DB.Save(&config.Account{
		Username: username,
		Password: passwordENC,
		AddedAt:  time.Now(),
		Settings: &config.Settings{
			Proxy: proxy,
		},
		CachedInsta: ig,
	})

	if err != nil {
		if err.Error() == "already exists" {
			http.Error(w, "An account with that username already exists.", http.StatusConflict)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf("Successfully added %s!", username)))
}

// APIAccountsHandler will return a list of all accounts marshalled into JSON
func APIAccountsHandler(w http.ResponseWriter, r *http.Request) {
	var accs = &[]config.Account{}
	err := config.DB.All(accs)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	b, err := json.Marshal(accs)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(b)
}

// APIDeleteAccountHandler is the http.Handler used to delete an account from the database
func APIDeleteAccountHandler(w http.ResponseWriter, r *http.Request) {
	username := chi.URLParam(r, "username") // /api/accounts/{username}
	acc := &config.Account{}

	err := config.DB.One("Username", username, acc)

	if err != nil {
		if err.Error() == "not found" {
			http.Error(w, "Account with the username not found", http.StatusNotFound)
			return
		}

		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = config.DB.DeleteStruct(acc)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write([]byte("Deleted " + username))
}
