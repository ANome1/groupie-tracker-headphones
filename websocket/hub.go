package websocket

import (
	"encoding/json"
	"log"
)

// BroadcastMessage represents a message to be sent to a specific room
type BroadcastMessage struct {
	RoomID  string
	Message []byte
}

// Hub maintains the set of active clients and broadcasts messages to the
// clients.
type Hub struct {
	// Registered clients.
	clients map[*Client]bool

	// Inbound messages from the clients.
	broadcast chan BroadcastMessage

	// Register requests from the clients.
	register chan *Client

	// Unregister requests from clients.
	unregister chan *Client
}

// NewHub creates a new Hub instance
func NewHub() *Hub {
	return &Hub{
		broadcast:  make(chan BroadcastMessage),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[*Client]bool),
	}
}

// Run starts the Hub's main loop
func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.clients[client] = true
		case client := <-h.unregister:
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)
			}
		case msg := <-h.broadcast:
			for client := range h.clients {
				if client.RoomID == msg.RoomID {
					select {
					case client.send <- msg.Message:
					default:
						close(client.send)
						delete(h.clients, client)
					}
				}
			}
		}
	}
}

// BroadcastToRoom sends a JSON message to all clients in a specific room
func (h *Hub) BroadcastToRoom(roomID string, v interface{}) {
	data, err := json.Marshal(v)
	if err != nil {
		log.Printf("Error marshalling broadcast message: %v", err)
		return
	}
	// Add newline delimiter for message parsing
	data = append(data, '\n')
	h.broadcast <- BroadcastMessage{
		RoomID:  roomID,
		Message: data,
	}
}
