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

	room, err := RoomService.GetRoomByCode(roomCode)
	if err != nil {
		http.Error(w, "Salle non trouvée", http.StatusNotFound)
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
		"Room":      room,
	}

	// Vérifier si une playlist a déjà été choisie par le host
	game := services.BlindTestMgr.GetGame(roomCode)
	if game != nil {
		if game.Config.PlaylistID != "" && game.Status != "finished" {
			// Une playlist a été définie, la pré-sélectionner
			data["PreSelectedPlaylist"] = game.Config.PlaylistID
		}
		// Passer la configuration actuelle
		data["CurrentNumRounds"] = game.Config.NumRounds
		data["CurrentTimePerRound"] = game.Config.TimePerRound
	} else {
		// Valeurs par défaut
		data["CurrentNumRounds"] = 5
		data["CurrentTimePerRound"] = 30
	}

	tmpl.Execute(w, data)
}
