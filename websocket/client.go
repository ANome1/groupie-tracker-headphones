package websocket

// TODO @Quoc Huy & @ilian: Structure Client
// TODO @Quoc Huy & @ilian: ReadPump, WritePump

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"groupie-tracker/models"
	"groupie-tracker/services"

	"github.com/gorilla/websocket"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 512
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true // Allow all origins for now
	},
}

// Client is a middleman between the websocket connection and the hub.
type Client struct {
	hub *Hub

	// The websocket connection.
	conn *websocket.Conn

	// Buffered channel of outbound messages.
	send chan []byte

	// Player ID
	PlayerID string
	// Room ID
	RoomID string
}

// ReadPump pumps messages from the websocket connection to the hub.
func (c *Client) ReadPump() {
	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
	}()
	c.conn.SetReadLimit(maxMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error { c.conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}

		c.handleMessage(message)
	}
}

func (c *Client) handleMessage(message []byte) {
	// Tenter de parser comme un message générique pour déterminer le type
	var genericMsg map[string]interface{}
	if err := json.Unmarshal(message, &genericMsg); err != nil {
		log.Printf("Error unmarshalling message: %v", err)
		return
	}

	// Vérifier si c'est un message blind test (type en minuscule)
	if msgType, ok := genericMsg["type"].(string); ok {
		// Messages blind test
		switch msgType {
		case "select_playlist", "submit_answer", "next_round":
			// Transmettre au BlindTestManager
			broadcastFunc := func(roomCode string, data interface{}) {
				// Convertir la map en JSON et envoyer avec newline delimiter
				jsonData, err := json.Marshal(data)
				if err != nil {
					log.Printf("Error marshalling blind test message: %v", err)
					return
				}
				// Ajouter une newline pour délimiter les messages
				jsonData = append(jsonData, '\n')
				c.hub.broadcast <- BroadcastMessage{
					RoomID:  roomCode,
					Message: jsonData,
				}
			}
			services.BlindTestMgr.HandleMessage(c.RoomID, c.PlayerID, message, broadcastFunc)
			return
		}
	}

	// Sinon, parser comme MessageIn (format Petit Bac)
	var msg models.MessageIn
	if err := json.Unmarshal(message, &msg); err != nil {
		log.Printf("Error unmarshalling MessageIn: %v", err)
		return
	}

	switch msg.Type {
	case "SUBMIT_ANSWERS":
		log.Printf("Received SUBMIT_ANSWERS raw data: %s", string(msg.Data))
		var payload models.AnswersPayload
		if err := json.Unmarshal(msg.Data, &payload); err != nil {
			log.Printf("Error unmarshalling payload: %v", err)
			return
		}
		log.Printf("Parsed AnswersPayload: %+v", payload)

		// Call service
		services.Manager.SubmitAnswers(c.RoomID, c.PlayerID, payload.Answers)

		// Check if we need to broadcast validation phase
		game := services.Manager.GetGame(c.RoomID)
		if game != nil {
			game.RLock()
			if game.State == models.GameStateVoting {
				// Get the current round to send responses
				if len(game.Rounds) > 0 {
					currentRound := game.Rounds[len(game.Rounds)-1]
					log.Printf("Broadcasting VALIDATION_PHASE to room %s. Responses count: %d", c.RoomID, len(currentRound.Responses))

					// Debug: Print responses
					for pid, resp := range currentRound.Responses {
						log.Printf("Player %s answers: %v", pid, resp.Answers)
					}

					if len(currentRound.Responses) == 0 {
						log.Printf("WARNING: Broadcasting VALIDATION_PHASE with EMPTY responses!")
					}

					// Include player names and host ID
					validationData := map[string]interface{}{
						"Round":       currentRound,
						"PlayerNames": game.PlayerNames,
						"HostID":      game.HostID,
					}

					c.hub.BroadcastToRoom(c.RoomID, models.MessageOut{
						Type: "VALIDATION_PHASE",
						Data: validationData,
					})
				}
			}
			game.RUnlock()
		}

	case "NEXT_ROUND":
		game := services.Manager.GetGame(c.RoomID)
		if game != nil && game.HostID == c.PlayerID {
			update, gameOver := services.Manager.NextRound(c.RoomID)
			if gameOver {
				c.hub.BroadcastToRoom(c.RoomID, models.MessageOut{
					Type: "GAME_OVER",
					Data: game.Scores,
				})
			} else if update != nil {
				c.hub.BroadcastToRoom(c.RoomID, models.MessageOut{
					Type: "NEW_ROUND",
					Data: update,
				})
			}
		}

	case "SUBMIT_VOTE":
		var vote models.Vote
		if err := json.Unmarshal(msg.Data, &vote); err != nil {
			log.Printf("Error unmarshalling vote: %v", err)
			return
		}
		vote.VoterID = c.PlayerID

		// Call service to record vote
		services.Manager.SubmitVote(c.RoomID, vote)

		// Calculate and broadcast scores (real-time update)
		scores := services.Manager.CalculateScores(c.RoomID)
		if scores != nil {
			c.hub.BroadcastToRoom(c.RoomID, models.MessageOut{
				Type: "SCORES_UPDATE",
				Data: scores,
			})
		}
	}
}

// WritePump pumps messages from the hub to the websocket connection.
func (c *Client) WritePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()
	for {
		select {
		case message, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// The hub closed the channel.
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)

			// Add queued chat messages to the current websocket message.
			n := len(c.send)
			for i := 0; i < n; i++ {
				w.Write(<-c.send)
			}

			if err := w.Close(); err != nil {
				return
			}
		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

// ServeWs handles websocket requests from the peer.
func ServeWs(hub *Hub, w http.ResponseWriter, r *http.Request) {
	// Extract roomID and playerID from query params
	roomID := r.URL.Query().Get("room")
	playerID := r.URL.Query().Get("player")

	if roomID == "" || playerID == "" {
		http.Error(w, "Missing room or player ID", http.StatusBadRequest)
		return
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	client := &Client{
		hub:      hub,
		conn:     conn,
		send:     make(chan []byte, 256),
		RoomID:   roomID,
		PlayerID: playerID,
	}
	client.hub.register <- client

	// Allow collection of memory referenced by the caller by doing all work in
	// new goroutines.
	go client.WritePump()
	go client.ReadPump()
}
