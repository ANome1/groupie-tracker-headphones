package websocket

// RESPONSABLE: @Quoc Huy & @ilian
// Client WebSocket individuel

// TODO @Quoc Huy & @ilian: Créer la structure Client avec:
// - hub *Hub
// - conn *websocket.Conn
// - send chan []byte
// - roomID int
// - userID int

// TODO @Quoc Huy & @ilian: Méthodes ReadPump() et WritePump()
// - ReadPump: Lire les messages du client et les traiter
// - WritePump: Envoyer les messages au client
