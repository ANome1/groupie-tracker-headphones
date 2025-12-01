package websocket

// RESPONSABLE: @Quoc Huy & @ilian
// Hub WebSocket central pour gérer toutes les connexions temps réel
// IMPORTANT: Ne pas utiliser github.com/gorilla/websocket

import "sync"

// Hub - Gère toutes les connexions WebSocket actives
type Hub struct {
	clients    map[*Client]bool         // Clients connectés
	rooms      map[int]map[*Client]bool // Clients par salle (roomID -> clients)
	broadcast  chan Message             // Canal pour diffuser des messages
	register   chan *Client             // Canal pour enregistrer un client
	unregister chan *Client             // Canal pour désenregistrer un client
	mu         sync.RWMutex             // Mutex pour la synchronisation
}

// NewHub - Crée un nouveau hub WebSocket
// TODO @Quoc Huy & @ilian:
func NewHub() *Hub {
	return &Hub{
		clients:    make(map[*Client]bool),
		rooms:      make(map[int]map[*Client]bool),
		broadcast:  make(chan Message),
		register:   make(chan *Client),
		unregister: make(chan *Client),
	}
}

// Run - Démarre le hub (boucle infinie)
// TODO @Quoc Huy & @ilian:
// - Écouter sur les canaux register, unregister, broadcast
// - Gérer l'ajout/suppression de clients
// - Diffuser les messages aux clients concernés
// IMPORTANT: Cette fonction doit tourner dans une goroutine
// Ex: go hub.Run()
func (h *Hub) Run() {
	// TODO: Implementation
	// for {
	//   select {
	//   case client := <-h.register:
	//     // Ajouter le client
	//   case client := <-h.unregister:
	//     // Retirer le client
	//   case message := <-h.broadcast:
	//     // Diffuser le message
	//   }
	// }
}

// BroadcastToRoom - Envoie un message à tous les clients d'une salle
// TODO @Quoc Huy & @ilian:
// - Utilisé pour les événements de jeu (nouvelle musique, nouveau round, scores, etc.)
func (h *Hub) BroadcastToRoom(roomID int, message Message) {
	// TODO: Implementation
}
