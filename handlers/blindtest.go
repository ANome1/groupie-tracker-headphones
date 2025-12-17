package handlers

import (
	"groupie-tracker/services"
	"groupie-tracker/utils"
	"html/template"
	"net/http"
	"sync"
)

type PlaylistOption struct {
	Name string
	ID   string
}

var (
	deezerGenres = []PlaylistOption{
		{"🎸 Rock", "152"},
		{"🎤 Pop", "132"},
		{"🎧 Rap US", "116"},
		{"🇫🇷 Rap FR", "162"},
		{"🎵 Indie", "172"},
	}
	genresMutex sync.RWMutex
)

func BlindTestHandler(w http.ResponseWriter, r *http.Request) {
	roomCode := r.URL.Query().Get("code")
	if roomCode == "" {
		http.Error(w, "Code de salle manquant", http.StatusBadRequest)
		return
	}

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

	game := services.BlindTestMgr.GetGame(roomCode)
	if game != nil {
		if game.Config.PlaylistID != "" && game.Status != "finished" {
			data["PreSelectedPlaylist"] = game.Config.PlaylistID
		}
		data["CurrentNumRounds"] = game.Config.NumRounds
		data["CurrentTimePerRound"] = game.Config.TimePerRound
	} else {
		data["CurrentNumRounds"] = 5
		data["CurrentTimePerRound"] = 30
	}

	tmpl.Execute(w, data)
}
