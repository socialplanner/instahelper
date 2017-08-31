package notifications

// READ https://www.jonathan-petitcolas.com/2015/01/27/playing-with-websockets-in-go.html

import (
	"github.com/socialplanner/instahelper/app/config"

	"encoding/json"
)

// Notification is a link to the underlying type of config.Notification
type Notification = config.Notification

// NewNotification will push a websocket message to all
func NewNotification(text, link string) error {
	n := &Notification{
		Text: text,
		Link: link,
	}

	err := config.DB.Save(n)

	if err != nil {
		return err
	}

	b, err := json.Marshal(n)

	if err != nil {
		return err
	}

	Hub.broadcast <- b
	return nil
}

// GetNotifications will return all notifications currently saved in the database
// - returns nil on error
func GetNotifications() *[]Notification {
	notifs := &[]Notification{}

	err := config.DB.All(notifs)

	if err != nil {
		return nil
	}

	return notifs
}
