package services

import (
	"encoding/json"
	"fmt"
	"groupie-tracker/models"
	"log"
	"strings"
	"sync"
	"time"
)

// BlindTestManager gère les sessions de blind test
type BlindTestManager struct {
	games map[string]*models.BlindTestGame
	mutex sync.RWMutex
}

// NewBlindTestManager crée un nouveau gestionnaire
func NewBlindTestManager() *BlindTestManager {
	return &BlindTestManager{
		games: make(map[string]*models.BlindTestGame),
	}
}

// CreateGame crée une nouvelle partie de blind test
func (m *BlindTestManager) CreateGame(roomCode string, hostID int, config models.BlindTestConfig) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	m.games[roomCode] = models.NewBlindTestGame(roomCode, hostID, config)
	log.Printf("Blind test créé pour room %s", roomCode)
}

// GetGame récupère une partie
func (m *BlindTestManager) GetGame(roomCode string) *models.BlindTestGame {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	return m.games[roomCode]
}

// HandleMessage traite les messages WebSocket pour le blind test
func (m *BlindTestManager) HandleMessage(roomCode, username string, message []byte, broadcastFunc func(string, interface{})) {
	var msg map[string]interface{}
	if err := json.Unmarshal(message, &msg); err != nil {
		log.Printf("Erreur parsing message: %v", err)
		return
	}

	msgType, ok := msg["type"].(string)
	if !ok {
		log.Printf("Message sans type")
		return
	}

	game := m.GetGame(roomCode)
	if game == nil {
		log.Printf("Jeu non trouvé pour room %s", roomCode)
		return
	}

	switch msgType {
	case "select_playlist":
		m.handleSelectPlaylist(game, msg, broadcastFunc)
	case "submit_answer":
		m.handleSubmitAnswer(game, username, msg, broadcastFunc)
	case "next_round":
		m.handleNextRound(game, broadcastFunc)
	default:
		log.Printf("Type de message inconnu: %s", msgType)
	}
}

// handleSelectPlaylist démarre le jeu avec la playlist choisie
func (m *BlindTestManager) handleSelectPlaylist(game *models.BlindTestGame, msg map[string]interface{}, broadcastFunc func(string, interface{})) {
	playlistID, ok := msg["playlist_id"].(string)
	if !ok {
		log.Printf("Playlist ID manquant")
		return
	}

	// Éviter de relancer la première manche si elle a déjà commencé
	if game.Status == "playing" && game.RoundNumber > 0 {
		log.Printf("Partie déjà en cours pour room %s, ignorant sélection de playlist", game.RoomCode)
		return
	}

	game.Config.PlaylistID = playlistID
	game.Status = "playing"

	// Démarrer la première manche
	m.startNewRound(game, broadcastFunc)
}

// startNewRound démarre une nouvelle manche
func (m *BlindTestManager) startNewRound(game *models.BlindTestGame, broadcastFunc func(string, interface{})) {
	game.RoundNumber++
	game.ResetRound()

	// Récupérer un track aléatoire depuis Deezer
	track, err := GetRandomTrackFromDeezerGenre(game.Config.PlaylistID)
	if err != nil {
		log.Printf("Erreur récupération track: %v", err)
		broadcastFunc(game.RoomCode, map[string]interface{}{
			"type":  "error",
			"error": "Impossible de charger la musique",
		})
		return
	}

	log.Printf("Track récupéré: ID=%d, Title=%s, Artist=%s, Preview=%s", track.ID, track.Title, track.Artist.Name, track.Preview)

	// Créer la manche
	game.CurrentRound = &models.BlindTestRound{
		RoundNumber:           game.RoundNumber,
		TrackID:               track.ID,
		TrackName:             track.Title,
		ArtistName:            track.Artist.Name,
		PreviewURL:            track.Preview,
		CoverImage:            track.Album.CoverMedium,
		Duration:              game.Config.TimePerRound,
		StartTime:             time.Now(),
		CorrectAnswersCount:   0,
		CorrectAnswersPlayers: []string{},
	}

	// Broadcast aux clients
	msg := map[string]interface{}{
		"type":         "round_start",
		"round_number": game.RoundNumber,
		"total_rounds": game.Config.NumRounds,
		"preview_url":  track.Preview,
		"cover_image":  track.Album.CoverMedium,
		"duration":     game.Config.TimePerRound,
		"track_name":   track.Title,
		"artist_name":  track.Artist.Name,
	}
	log.Printf("Broadcasting round_start message: %+v", msg)
	broadcastFunc(game.RoomCode, msg)

	// Timer automatique pour finir la manche
	go func() {
		time.Sleep(time.Duration(game.Config.TimePerRound) * time.Second)
		m.endRound(game, broadcastFunc)
		m.checkGameEnd(game, broadcastFunc)
	}()
}

// handleSubmitAnswer traite la réponse d'un joueur
func (m *BlindTestManager) handleSubmitAnswer(game *models.BlindTestGame, username string, msg map[string]interface{}, broadcastFunc func(string, interface{})) {
	if game.CurrentRound == nil {
		return
	}

	// Vérifier si le joueur a déjà répondu
	if game.AnsweredPlayers[username] {
		return
	}

	trackAnswer, _ := msg["track"].(string)
	artistAnswer, _ := msg["artist"].(string)

	// Vérifier si la réponse est correcte
	isCorrect := m.checkAnswer(trackAnswer, artistAnswer, game.CurrentRound)

	// Calculer le rang (position parmi les réponses correctes)
	rank := game.CurrentRound.CorrectAnswersCount + 1

	// Enregistrer la réponse avec le rang
	game.RecordAnswer(username, isCorrect, rank)

	log.Printf("Joueur %s a répondu: correct=%v, rang=%d", username, isCorrect, rank)
}

// checkAnswer vérifie si la réponse est correcte
func (m *BlindTestManager) checkAnswer(trackAnswer, artistAnswer string, round *models.BlindTestRound) bool {
	trackAnswer = strings.ToLower(strings.TrimSpace(trackAnswer))
	artistAnswer = strings.ToLower(strings.TrimSpace(artistAnswer))

	// Vérifier s'il y a au moins une réponse
	if trackAnswer == "" && artistAnswer == "" {
		return false
	}

	trackMatch := m.fuzzyMatch(trackAnswer, round.TrackName)
	artistMatch := m.fuzzyMatch(artistAnswer, round.ArtistName)

	// Au moins l'un des deux doit être correct
	// (Accepter soit le titre, soit l'artiste)
	return trackMatch || artistMatch
}

// fuzzyMatch compare deux chaînes avec une tolérance
func (m *BlindTestManager) fuzzyMatch(answer, correct string) bool {
	answer = strings.ToLower(strings.TrimSpace(answer))
	correct = strings.ToLower(strings.TrimSpace(correct))

	if answer == "" {
		return false
	}

	// Match exact
	if answer == correct {
		return true
	}

	// Minimum length check: réponse doit être au moins 3 caractères
	if len(answer) < 3 {
		return false
	}

	// Match sur mots complets (pour éviter "on" matching "Song")
	words := strings.Fields(correct)
	for _, word := range words {
		// Match exact sur un mot
		if word == answer {
			return true
		}
		// Match partiel si la réponse est au moins 60% du mot
		similarity := levenshteinSimilarity(answer, word)
		if similarity >= 0.6 && len(answer) >= 3 {
			return true
		}
	}

	// Match si la réponse contient au moins 60% de la réponse correcte
	if len(correct) >= 3 {
		similarity := levenshteinSimilarity(answer, correct)
		if similarity >= 0.6 {
			return true
		}
	}

	return false
}

// levenshteinSimilarity calcule la similitude entre deux chaînes (0 à 1)
func levenshteinSimilarity(s1, s2 string) float64 {
	dist := levenshteinDistance(s1, s2)
	maxLen := len(s1)
	if len(s2) > maxLen {
		maxLen = len(s2)
	}
	if maxLen == 0 {
		return 1.0
	}
	return 1.0 - float64(dist)/float64(maxLen)
}

// levenshteinDistance calcule la distance d'édition entre deux chaînes
func levenshteinDistance(s1, s2 string) int {
	if len(s1) < len(s2) {
		return levenshteinDistance(s2, s1)
	}

	if len(s2) == 0 {
		return len(s1)
	}

	prev := make([]int, len(s2)+1)
	for i := 0; i <= len(s2); i++ {
		prev[i] = i
	}

	for i := 1; i <= len(s1); i++ {
		curr := make([]int, len(s2)+1)
		curr[0] = i

		for j := 1; j <= len(s2); j++ {
			cost := 0
			if s1[i-1] != s2[j-1] {
				cost = 1
			}

			curr[j] = min(
				curr[j-1]+1, // insertion
				min(prev[j]+1, // deletion
					prev[j-1]+cost)) // substitution
		}
		prev = curr
	}

	return prev[len(s2)]
}

// min retourne le minimum de deux entiers
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// endRound termine la manche et envoie les résultats
func (m *BlindTestManager) endRound(game *models.BlindTestGame, broadcastFunc func(string, interface{})) {
	if game.CurrentRound == nil {
		return
	}

	// Calculer les points de cette manche
	roundScores := m.calculateRoundScores(game)

	// Ajouter à l'historique
	game.RoundHistory = append(game.RoundHistory, *game.CurrentRound)

	// Broadcast les résultats avec les points de la manche
	broadcastFunc(game.RoomCode, map[string]interface{}{
		"type":            "round_end",
		"correct_track":   game.CurrentRound.TrackName,
		"correct_artist":  game.CurrentRound.ArtistName,
		"cover_image":     game.CurrentRound.CoverImage,
		"scores":          game.GetScoreboard(),
		"round_scores":    roundScores,
		"correct_players": game.CurrentRound.CorrectAnswersPlayers,
		"host_id":         game.HostID,
	})

	// Mettre à jour le scoreboard
	broadcastFunc(game.RoomCode, map[string]interface{}{
		"type":   "scoreboard_update",
		"scores": game.GetScoreboard(),
	})
}

// calculateRoundScores retourne les points gagnés par chaque joueur cette manche
func (m *BlindTestManager) calculateRoundScores(game *models.BlindTestGame) map[string]int {
	roundScores := make(map[string]int)

	if game.CurrentRound == nil {
		return roundScores
	}

	pointsByRank := []int{100, 80, 60, 40, 20}

	for idx, username := range game.CurrentRound.CorrectAnswersPlayers {
		points := 0
		if idx < len(pointsByRank) {
			points = pointsByRank[idx]
		}
		roundScores[username] = points
	}

	return roundScores
}

// Vérifier si c'est la fin du jeu après les résultats
func (m *BlindTestManager) checkGameEnd(game *models.BlindTestGame, broadcastFunc func(string, interface{})) {
	if game.RoundNumber >= game.Config.NumRounds {
		game.Status = "finished"
		broadcastFunc(game.RoomCode, map[string]interface{}{
			"type":   "game_end",
			"scores": game.GetScoreboard(),
		})
	}
}

// handleNextRound démarre la manche suivante
func (m *BlindTestManager) handleNextRound(game *models.BlindTestGame, broadcastFunc func(string, interface{})) {
	if game.Status == "finished" {
		return
	}

	if game.RoundNumber >= game.Config.NumRounds {
		game.Status = "finished"
		return
	}

	m.startNewRound(game, broadcastFunc)
}

// RemoveGame supprime une partie
func (m *BlindTestManager) RemoveGame(roomCode string) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	delete(m.games, roomCode)
	fmt.Printf("Blind test supprimé pour room %s\n", roomCode)
}
