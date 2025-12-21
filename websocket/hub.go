package websocket

import (
	"encoding/json"
	"log"
)

type BroadcastMessage struct {
	RoomID  string
	Message []byte
}

// Hub gère les connexions WebSocket et diffuse les messages par salle
type Hub struct {
	clients    map[*Client]bool
	broadcast  chan BroadcastMessage
	register   chan *Client
	unregister chan *Client
}

func NewHub() *Hub {
	return &Hub{
		broadcast:  make(chan BroadcastMessage),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[*Client]bool),
	}
}

// Boucle centrale: enregistre/désenregistre les clients et diffuse les messages aux salles
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

// Diffuse un message JSON à tous les clients d'une salle (avec délimiteur newline)
func (h *Hub) BroadcastToRoom(roomID string, v interface{}) {
	data, err := json.Marshal(v)
	if err != nil {
		log.Printf("Error marshalling broadcast message: %v", err)
		return
	}
	data = append(data, '\n')
	h.broadcast <- BroadcastMessage{
		RoomID:  roomID,
		Message: data,
	}
}

// Diffuse un message brut à tous les clients d'une salle
func (h *Hub) Broadcast(msg BroadcastMessage) {
	h.broadcast <- msg
}
