package websocket

// RESPONSABLE: @Quoc Huy & @ilian
// Hub central pour gérer les connexions WebSocket

// IMPORTANT: Ne pas utiliser github.com/gorilla/websocket
// Utiliser uniquement le package standard Go

// TODO @Quoc Huy & @ilian: Créer la structure Hub avec:
// - clients map[*Client]bool (clients connectés)
// - rooms map[int]map[*Client]bool (roomID -> clients)
// - broadcast chan []byte (canal pour broadcast)
// - register chan *Client (canal pour enregistrer)
// - unregister chan *Client (canal pour désenregistrer)

// TODO @Quoc Huy & @ilian: Fonction Run()
// - Boucle infinie qui écoute les canaux
// - Gérer l'enregistrement et le désenregistrement des clients
// - Gérer le broadcast des messages
