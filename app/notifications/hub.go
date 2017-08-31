package notifications

// HubType is the type for a websocket hub
type HubType struct {
	clients      map[*Client]bool
	broadcast    chan []byte
	addClient    chan *Client
	removeClient chan *Client
}

// Start runs forever as a goroutine
func (hub *HubType) Start() {
	for {
		// one of these fires when a channel
		// receives data
		select {
		case conn := <-hub.addClient:
			// add a new client
			hub.clients[conn] = true
		case conn := <-hub.removeClient:
			// remove a client
			if _, ok := hub.clients[conn]; ok {
				delete(hub.clients, conn)
				close(conn.send)
			}
		case message := <-hub.broadcast:
			// broadcast a message to all clients
			for conn := range hub.clients {
				select {
				case conn.send <- message:
				default:
					close(conn.send)
					delete(hub.clients, conn)
				}
			}
		}
	}
}

// Hub is a websocket hub
//
// To send message do Hub.broadcast <- []byte("MESSAGE")
var Hub = HubType{
	broadcast:    make(chan []byte),
	addClient:    make(chan *Client),
	removeClient: make(chan *Client),
	clients:      make(map[*Client]bool),
}
