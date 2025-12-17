package handlers

import (
	"groupie-tracker/services"
	"groupie-tracker/utils"
	"html/template"
	"net/http"
	"sync"
)

// Playlist structure to maintain order
type PlaylistOption struct {
	Name string
	ID   string
}

// Genres Deezer disponibles (ordered list)
var (
	deezerGenres = []PlaylistOption{
		{"🎸 Rock", "152"},    // Rock
		{"🎤 Pop", "132"},     // Pop
		{"🎧 Rap US", "116"},  // Rap US
		{"🇫🇷 Rap FR", "162"}, // Rap Français
		{"🎵 Indie", "172"},   // Indie
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

	// Vérifier si une playlist a déjà été choisie par le host
	game := services.BlindTestMgr.GetGame(roomCode)
	if game != nil && game.Config.PlaylistID != "" {
		// Une playlist a été définie, la pré-sélectionner
		data["PreSelectedPlaylist"] = game.Config.PlaylistID
	}

	tmpl.Execute(w, data)
}
