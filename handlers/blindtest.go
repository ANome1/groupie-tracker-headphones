package handlers

import (
	"groupie-tracker/utils"
	"html/template"
	"net/http"
	"sync"
)

// Genres Deezer disponibles
var (
	deezerGenres = map[string]string{
		"🎸 Rock": "152", // Rock genre ID
		"🎤 Pop":  "132", // Pop genre ID
		"🎧 Rap":  "116", // Rap/Hip-Hop genre ID
	}
	genresMutex sync.RWMutex
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

	genresMutex.RLock()
	defer genresMutex.RUnlock()

	data := map[string]interface{}{
		"Playlists": deezerGenres,
		"RoomCode":  roomCode,
		"User":      user,
	}

	tmpl.Execute(w, data)
}
