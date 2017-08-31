package notifications

import "github.com/gorilla/websocket"

// Client is a websocket client
type Client struct {
	ws *websocket.Conn
	// Hub passes broadcast messages to this channel
	send chan []byte
}

// Hub broadcasts a new message and this fires
func (c *Client) write() {
	// make sure to close the connection incase the loop exits
	defer func() {
		Hub.removeClient <- c
		c.ws.Close()
	}()

	for {
		select {
		case message, ok := <-c.send:
			if !ok {
				c.ws.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			c.ws.WriteMessage(websocket.TextMessage, message)
		}
	}
}

// // New message received so pass it to the Hub
// func (c *Client) read() {
// 	defer func() {
// 		Hub.removeClient <- c
// 		c.ws.Close()
// 	}()

// 	for {
// 		_, message, err := c.ws.ReadMessage()
// 		if err != nil {
// 			Hub.removeClient <- c
// 			c.ws.Close()
// 			break
// 		}

// 		Hub.broadcast <- message
// 	}
// }
