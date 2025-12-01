package websocket

// RESPONSABLE: @Quoc Huy & @ilian
// Client WebSocket individuel

// Client - Représente une connexion WebSocket d'un joueur
type Client struct {
	hub    *Hub        // Référence au hub
	conn   interface{} // Connexion WebSocket (remplacer par le type approprié)
	send   chan []byte // Canal pour envoyer des messages
	UserID int         // ID de l'utilisateur connecté
	RoomID int         // ID de la salle où se trouve le joueur
}

// NewClient - Crée un nouveau client WebSocket
// TODO @Quoc Huy & @ilian:
func NewClient(hub *Hub, conn interface{}, userID, roomID int) *Client {
	return &Client{
		hub:    hub,
		conn:   conn,
		send:   make(chan []byte, 256),
		UserID: userID,
		RoomID: roomID,
	}
}

// ReadPump - Lit les messages entrants du client
// TODO @Quoc Huy & @ilian:
// - Boucle infinie qui lit les messages WebSocket
// - Parser les messages JSON (voir messages.go)
// - Router les messages selon leur type:
//   - Blind Test: soumission de réponse, demande de round suivant
//   - Petit Bac: soumission de réponses, votes de validation
//
// - Gérer les erreurs de connexion
// IMPORTANT: Lancer dans une goroutine
// Ex: go client.ReadPump()
func (c *Client) ReadPump() {
	// TODO: Implementation
	// defer func() {
	//   c.hub.unregister <- c
	//   c.conn.Close()
	// }()
	// for {
	//   // Lire le message
	//   // Parser le JSON
	//   // Router selon le type
	// }
}

// WritePump - Écrit les messages sortants vers le client
// TODO @Quoc Huy & @ilian:
// - Boucle infinie qui envoie les messages du canal 'send'
// - Gérer les erreurs d'écriture
// - Implémenter un ping/pong pour détecter les déconnexions
// IMPORTANT: Lancer dans une goroutine
// Ex: go client.WritePump()
func (c *Client) WritePump() {
	// TODO: Implementation
	// for {
	//   select {
	//   case message := <-c.send:
	//     // Envoyer le message
	//   }
	// }
}
