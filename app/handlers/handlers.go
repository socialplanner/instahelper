package handlers

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/socialplanner/instahelper/app/assets"
	"github.com/socialplanner/instahelper/app/config"

	"github.com/sirupsen/logrus"
)

// Pages is a collection of all pages of instahelper
var Pages = map[string]Page{}

func init() {
	// Pages is a collection of all pages of instahelper
	Pages = map[string]Page{
		// Main Dashboard
		"dashboard": {
			ID:       1,
			Name:     "Dashboard",
			Link:     "/",
			Icon:     "dashboard",
			Template: newTemplate("base.html", "dashboard.html"),
			Handler:  DashboardHandler,
		},

		"register": {
			ID:       2,
			Name:     "Add Account",
			Link:     "/register",
			Icon:     "person_add",
			Template: newTemplate("base.html", "register.html"),
			Handler:  RegisterHandler,
		},

		"accounts": {
			ID:       3,
			Name:     "Accounts",
			Link:     "/accounts",
			Icon:     "people",
			Template: newTemplate("base.html", "accounts.html"),
			Handler:  AccountsHandler,
		},

		"account": {
			ID:       4,
			Name:     "Account Settings",
			Link:     "/accounts/{username}",
			Icon:     "settings",
			Template: newTemplate("base.html", "account.html"),
			Unlisted: true,
			Handler:  AccSettingsHandler,
		},

		"settings": {
			ID:       5,
			Name:     "Settings",
			Link:     "/settings",
			Icon:     "settings",
			Template: newTemplate("base.html", "settings.html"),
			Handler:  SettingsHandler,
		},

		"update": {
			ID:       6,
			Name:     "Update",
			Link:     "/update",
			Icon:     "get_app",
			Template: newTemplate("base.html", "update.html"),
			Handler:  UpdateHandler,
		},
	}
}

// DashboardHandler is the handler for the main dashboard of Instahelper
func DashboardHandler(w http.ResponseWriter, r *http.Request) {
	err := Template("dashboard").Execute(w)

	if err != nil {
		logrus.Error(err)
		Error(w, err)
	}
}

// RegisterHandler is the handler for the "Add Account" page of Instahelper
func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	err := Template("register").Execute(w)

	if err != nil {
		logrus.Error(err)
		Error(w, err)
	}
}

// AccountsHandler is the handler for the "Accounts" page of Instahelper
func AccountsHandler(w http.ResponseWriter, r *http.Request) {
	var accs = &[]config.Account{}
	err := config.DB.All(accs)

	if err != nil {
		logrus.Error(err)
		Error(w, err)
		return
	}

	err = Template("accounts").Execute(w, map[string]interface{}{
		"Accounts": accs,
	})

	if err != nil {
		logrus.Error(err)
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
		http.Error(w, "asset not found", http.StatusNotFound)
		return
	}

	w.Write(b)
}

// Error will display the error and promt the user to report it
// Do not use this for api endpoints. Only for user facing pages.
func Error(w http.ResponseWriter, e error) {
	w.WriteHeader(http.StatusInternalServerError)
	w.Write([]byte(
		fmt.Sprintf("Error: %s\nIt would be appreciated if you could create an issue at https://github.com/socialplanner/instahelper/issues/new", e.Error()),
	))
}
