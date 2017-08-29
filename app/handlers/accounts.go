package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

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
		w.Write([]byte("Invalid Form Input"))
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	acc := &config.Account{}

	config.DB.One("Username", username, acc)

	if acc.Username == username {
		w.WriteHeader(http.StatusConflict)
		w.Write([]byte("An account with that username already exists."))
		return
	}

	// Encrypt password
	passwordENC, err := insta.Encrypt(password)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	ig, err := insta.CachedInsta(username, password, proxy)

	if err != nil {
		// Bad username/password combo or captcha
		w.WriteHeader(http.StatusBadGateway)
		w.Write([]byte(err.Error()))
		return
	}

	err = config.DB.Save(&config.Account{
		Username:    username,
		Password:    passwordENC,
		AddedAt:     time.Now(),
		Settings:    &config.Settings{},
		CachedInsta: ig,
	})

	if err != nil {
		if err.Error() == "already exists" {
			w.WriteHeader(http.StatusConflict)
			w.Write([]byte("An account with that username already exists."))
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
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
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	b, err := json.Marshal(accs)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(b)
}
