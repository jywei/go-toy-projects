package chat

import (
	"net/http"

	"github.com/gorilla/websocket"
)

// WSUpgrader is used upgrade the protocol to allow websockets
var WSUpgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

// Conn represents a websocket connection and send message to them (clients)
// three methods: reading to the hub, writing to the hub, writing to the client
type Conn struct {
	WS   *websocket.Conn
	Send chan string
}

// SendToHub sends any message from our websocket connection to our hub
func (conn *Conn) SendToHub() {
	defer conn.WS.Close()
	for {
		_, msg, err := conn.WS.ReadMessage()
		if err != nil {
			// user has disconnected - they probably just refreshed their
			// browser, so just return
			return
		}
		DefaultHub.Echo <- string(msg)
	}
}

// ReceiveFromHub sends messages from our hub to our websocket connection
func (conn *Conn) ReceiveFromHub() {
	defer conn.WS.Close()
	for {
		conn.Write(<-conn.Send)
	}
}

// Write writes to clients
func (conn *Conn) Write(msg string) error {
	return conn.WS.WriteMessage(websocket.TextMessage, []byte(msg))
}

// WSHandler handles the HTTP req
func WSHandler(w http.ResponseWriter, r *http.Request) {
	// upgrade the connection
	ws, err := WSUpgrader.Upgrade(w, r, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// create a new connection
	conn := &Conn{
		Send: make(chan string),
		WS:   ws,
	}
	// add the connection to the hub
	DefaultHub.Join <- conn

	// send messages to the hub
	go conn.SendToHub()
	// and receiving from the hub at the meantime
	conn.ReceiveFromHub()
}
