package handlers

// RESPONSABLE: @Quoc Huy & @ilian
// Handlers pour les connexions WebSocket (communication temps réel)

import "net/http"

// WebSocketHandler - Gère les connexions WebSocket
// TODO @Quoc Huy & @ilian:
// - Upgrader la connexion HTTP en WebSocket
// - Vérifier l'authentification du joueur
// - Récupérer le roomID depuis les paramètres
// - Créer un nouveau client WebSocket
// - L'enregistrer dans le hub
// - Démarrer les goroutines de lecture/écriture
//
// IMPORTANT: Ne pas utiliser github.com/gorilla/websocket
// Utiliser le package net/http standard ou implémenter WebSocket manuellement
func WebSocketHandler(w http.ResponseWriter, r *http.Request) {
	// TODO: Implementation
	// Exemple de structure:
	// 1. Upgrader la connexion
	// 2. client := websocket.NewClient(hub, conn, userID, roomID)
	// 3. hub.Register <- client
	// 4. go client.ReadPump()
	// 5. go client.WritePump()
}
