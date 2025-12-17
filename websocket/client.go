package websocket

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
	writeWait      = 10 * time.Second
	pongWait       = 60 * time.Second
	pingPeriod     = (pongWait * 9) / 10
	maxMessageSize = 512
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type Client struct {
	hub      *Hub
	conn     *websocket.Conn
	send     chan []byte
	PlayerID string
	RoomID   string
}

// Boucle de lecture WebSocket: traite les messages reçus du client et les achemine vers les services
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
		case "join_game", "select_playlist", "submit_answer", "next_round":
			// Tenter d'extraire le username du message pour les submit_answer
			var userMsg map[string]interface{}
			json.Unmarshal(message, &userMsg)

			username := c.PlayerID // Par défaut, utiliser l'ID
			if msgType == "submit_answer" {
				if user, ok := userMsg["username"].(string); ok {
					username = user // Utiliser le pseudo si disponible
				}
			}

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

			sendToClient := func(v interface{}) {
				data, err := json.Marshal(v)
				if err == nil {
					data = append(data, '\n')
					c.send <- data
				}
			}

			services.BlindTestMgr.HandleMessage(c.RoomID, username, message, broadcastFunc, sendToClient)
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
	case "START_GAME":
		// Parse config from payload
		var configPayload struct {
			Config models.PetitBacConfig `json:"config"`
		}
		if err := json.Unmarshal(msg.Data, &configPayload); err == nil {
			// Update game config
			services.Manager.UpdateGameConfig(c.RoomID, configPayload.Config)
		}

		// Start the first round
		roundUpdate := services.Manager.StartRound(c.RoomID)
		if roundUpdate != nil {
			c.hub.BroadcastToRoom(c.RoomID, models.MessageOut{
				Type: "ROUND_START",
				Data: roundUpdate,
			})
		}

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

	case "END_ROUND":
		// Broadcaster les résultats détaillés de la manche
		game := services.Manager.GetGame(c.RoomID)
		if game != nil {
			roundResults := services.Manager.CalculateRoundResults(c.RoomID)
			if roundResults != nil {
				// Update actual game scores
				game.Lock()
				for pID, total := range roundResults.TotalScores {
					game.Scores[pID] = total
				}
				game.Unlock()

				// Broadcast round results to all players
				c.hub.BroadcastToRoom(c.RoomID, models.MessageOut{
					Type: "ROUND_RESULTS",
					Data: roundResults,
				})
			}
		}

	case "NEXT_ROUND":
		game := services.Manager.GetGame(c.RoomID)
		if game != nil && game.HostID == c.PlayerID {
			// Check if someone completed the round or it's the host action
			update, gameOver := services.Manager.NextRound(c.RoomID)
			if gameOver {
				// Game finished - broadcast final scores
				c.hub.BroadcastToRoom(c.RoomID, models.MessageOut{
					Type: "GAME_OVER",
					Data: map[string]interface{}{
						"scores": game.Scores,
						"hostID": game.HostID,
					},
				})
			} else if update != nil {
				// New round started
				c.hub.BroadcastToRoom(c.RoomID, models.MessageOut{
					Type: "NEW_ROUND",
					Data: update,
				})
			}
		}
	}
}

// Boucle d'écriture WebSocket: envoie les messages du hub vers le client
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
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)

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

// Établit la connexion WebSocket et lance les pompes de lecture/écriture
func ServeWs(hub *Hub, w http.ResponseWriter, r *http.Request) {
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

	go client.WritePump()
	go client.ReadPump()
}
