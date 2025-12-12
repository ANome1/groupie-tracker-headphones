package websocket

import (
"encoding/json"
"log"
"net/http"
"time"

"github.com/gorilla/websocket"
"groupie-tracker/models"
"groupie-tracker/services"
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
var msg models.MessageIn
if err := json.Unmarshal(message, &msg); err != nil {
log.Printf("Error unmarshalling message: %v", err)
return
}

switch msg.Type {
case "SUBMIT_ANSWERS":
var payload models.AnswersPayload
if err := json.Unmarshal(msg.Data, &payload); err != nil {
log.Printf("Error unmarshalling payload: %v", err)
return
}

// Call service
services.Manager.SubmitAnswers(c.RoomID, c.PlayerID, payload.Answers)

// Check if we need to broadcast validation phase
game := services.Manager.GetGame(c.RoomID)
if game != nil && game.State == models.GameStateVoting {
// Get the current round to send responses
if len(game.Rounds) > 0 {
currentRound := game.Rounds[len(game.Rounds)-1]
c.hub.BroadcastToRoom(c.RoomID, models.MessageOut{
Type: "VALIDATION_PHASE",
Data: currentRound,
})
}
}

case "SUBMIT_VOTE":
// TODO: Implement vote submission in service
// services.Manager.SubmitVote(...)
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
