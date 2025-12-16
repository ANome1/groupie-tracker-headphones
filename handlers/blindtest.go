package handlers

import (
	"groupie-tracker/utils"
	"html/template"
	"net/http"
	"sync"
)

// Playlists prédéfinies Deezer (ID de playlists publiques)
var (
	deezerPlaylists = map[string]string{
		"🎸 Rock": "1116188121",
		"🎤 Pop":  "1266970711",
		"🎧 Rap":  "1282483245",
	}
	playlistsMutex sync.RWMutex
)

// BlindTestHandler renders the blind test game page
func BlindTestHandler(w http.ResponseWriter, r *http.Request) {
	// Récupérer le code de la room
	roomCode := r.URL.Query().Get("code")
	if roomCode == "" {
		http.Error(w, "Code de salle manquant", http.StatusBadRequest)
		return
	}

	// Récupérer l'utilisateur
	userID, err := utils.GetUserIDFromCookie(r)
	if err != nil {
		http.Error(w, "Non authentifié", http.StatusUnauthorized)
		return
	}

	user, err := AuthService.GetUserByID(userID)
	if err != nil {
		http.Error(w, "Utilisateur non trouvé", http.StatusNotFound)
		return
	}

	tmpl, err := template.ParseFiles(
		"templates/games/blindtest.html",
		"templates/components/header.html",
		"templates/components/footer.html",
	)
	if err != nil {
		http.Error(w, "Error loading template: "+err.Error(), http.StatusInternalServerError)
		return
	}

	playlistsMutex.RLock()
	defer playlistsMutex.RUnlock()

	data := map[string]interface{}{
		"Playlists": deezerPlaylists,
		"RoomCode":  roomCode,
		"User":      user,
	}

	tmpl.Execute(w, data)
}
