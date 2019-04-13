package chat

// DefaultHub is our default hub
var DefaultHub = NewHub()

// Hub is the data structure we use to keep track of connections
type Hub struct {
	// Join channel: signal there is a new client attempting tpo connect
	Join chan *Conn
	// Conns: maintain currently active users
	Conns map[*Conn]bool
	// Echo channel: receive messages
	Echo chan string
}

// NewHub creates a new default hub.
func NewHub() *Hub {
	return &Hub{
		Join:  make(chan *Conn),
		Conns: make(map[*Conn]bool),
		Echo:  make(chan string),
	}
}

// Start starts our hub
func (hub *Hub) Start() {
	for {
		// select is a multiplexer for channel
		// it will wait for one of its cases to run
		select {
		// Join: add the connection to the hub
		case conn := <-hub.Join:
			DefaultHub.Conns[conn] = true
		// Echo: when receiving the message, send to all connected channels
		case msg := <-hub.Echo:
			for conn := range hub.Conns {
				conn.Send <- msg
			}
		}
	}
}
