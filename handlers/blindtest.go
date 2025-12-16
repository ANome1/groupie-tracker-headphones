package handlers

import (
	"html/template"
	"net/http"
	"sync"
)

// Playlists prédéfinies Deezer (ID de playlists publiques)
var (
	deezerPlaylists = map[string]string{
		"Top France":        "1306931615",
		"Top Monde":         "1266970711",
		"Hits Années 2000":  "1282483245",
		"Variété Française": "1109890031",
		"Rock Classics":     "1116188121",
	}
	playlistsMutex sync.RWMutex
)

// BlindTestHandler renders the blind test game page
func BlindTestHandler(w http.ResponseWriter, r *http.Request) {
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
	}

	tmpl.Execute(w, data)
}
