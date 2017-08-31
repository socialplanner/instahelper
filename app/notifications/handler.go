package notifications

import (
	"net/http"

	"github.com/gorilla/websocket"
)

// Upgrader is used to upgrade requests to websocket requests
var Upgrader = websocket.Upgrader{}

// WSHandler is the web socket handler
func WSHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := Upgrader.Upgrade(w, r, nil)

	if err != nil {
		http.NotFound(w, r)
		return
	}

	client := &Client{
		ws:   conn,
		send: make(chan []byte),
	}

	Hub.addClient <- client

	go client.write()
}
