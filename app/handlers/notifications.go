package handlers

import (
	"net/http"

	"github.com/socialplanner/instahelper/app/config"
)

// APIDeleteNotificationsHandler is a handler to delete all notifications from the database
func APIDeleteNotificationsHandler(w http.ResponseWriter, r *http.Request) {
	notifs := &[]config.Notification{}

	err := config.DB.All(notifs)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	for _, n := range *notifs {
		err = config.DB.DeleteStruct(&n)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	w.Write([]byte("Deleted All Notifications"))
}
