package handlers

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/socialplanner/instahelper/app/assets"
	"github.com/socialplanner/instahelper/app/config"

	l "github.com/socialplanner/instahelper/app/log"
)

var (
	log = l.Log
)

// Pages is a collection of all pages of instahelper
var Pages = map[string]Page{
	// Main Dashboard
	"dashboard": {
		ID:       1,
		Name:     "Dashboard",
		Link:     "/",
		Icon:     "dashboard",
		Template: newTemplate("base.html", "dashboard.html"),
	},

	"register": {
		ID:       2,
		Name:     "Add Account",
		Link:     "/register",
		Icon:     "person_add",
		Template: newTemplate("base.html", "register.html"),
	},

	"accounts": {
		ID:       3,
		Name:     "Accounts",
		Link:     "/accounts",
		Icon:     "people",
		Template: newTemplate("base.html", "accounts.html"),
	},
}

// Handler returns the http.Handler for the corresponding page
func (p *Page) Handler() func(http.ResponseWriter, *http.Request) {
	switch p.Name {
	case "Dashboard":
		return DashboardHandler
	case "Add Account":
		return RegisterHandler
	case "Accounts":
		return AccountsHandler
	}

	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		http.Redirect(w, r, "/", http.StatusPermanentRedirect)
	}
}

// DashboardHandler is the handler for the main dashboard of Instahelper
func DashboardHandler(w http.ResponseWriter, r *http.Request) {
	err := Template("dashboard").Execute(w)

	if err != nil {
		log.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		Error(w, err)
	}
}

// RegisterHandler is the handler for the "Add Account" page of Instahelper
func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	err := Template("register").Execute(w)

	if err != nil {
		log.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		Error(w, err)
	}
}

// AccountsHandler is the handler for the "Accounts" page of Instahelper
func AccountsHandler(w http.ResponseWriter, r *http.Request) {
	var accs = &[]config.Account{}
	err := config.DB.All(accs)

	if err != nil {
		log.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		Error(w, err)
	}

	err = Template("accounts").Execute(w, map[string]interface{}{
		"Accounts": accs,
	})

	if err != nil {
		log.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		Error(w, err)
	}
}

// AssetsHandler is the handler for all assets packed into the binary with go-bindata
func AssetsHandler(w http.ResponseWriter, r *http.Request) {
	url := r.URL.String()

	url = strings.Replace(url, "/assets/", "", -1)

	b, err := assets.Asset(url)

	if err != nil {
		// asset not found
		w.WriteHeader(http.StatusNotFound)
		Error(w, err)
	}

	w.Write(b)
}

// Error will display the error and promt the user to report it
func Error(w http.ResponseWriter, e error) {
	w.Write([]byte(
		fmt.Sprintf("Error: %s\nIt would be appreciated if you could create an issue at https://github.com/socialplanner/instahelper/issues/new", e.Error()),
	))
}
