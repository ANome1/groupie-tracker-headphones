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

type BlindTestManager struct {
	games map[string]*models.BlindTestGame
	mutex sync.RWMutex
}

func NewBlindTestManager() *BlindTestManager {
	return &BlindTestManager{
		games: make(map[string]*models.BlindTestGame),
	}
}

func (m *BlindTestManager) CreateGame(roomCode string, hostID int, config models.BlindTestConfig) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	m.games[roomCode] = models.NewBlindTestGame(roomCode, hostID, config)
	log.Printf("Blind test created for room %s", roomCode)
}

func (m *BlindTestManager) GetGame(roomCode string) *models.BlindTestGame {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	return m.games[roomCode]
}

// Synchronisation WebSocket: Reçoit les messages depuis les clients et les traite selon le type
// Routes les messages join_game, select_playlist, submit_answer, next_round vers leurs handlers
func (m *BlindTestManager) HandleMessage(roomCode, username string, message []byte, broadcastFunc func(string, interface{}), sendToClientFunc func(interface{})) {
	var msg map[string]interface{}
	if err := json.Unmarshal(message, &msg); err != nil {
		log.Printf("Error parsing message: %v", err)
		return
	}

	msgType, ok := msg["type"].(string)
	if !ok {
		log.Printf("Message missing type")
		return
	}

	game := m.GetGame(roomCode)
	if game == nil {
		log.Printf("Game not found for room %s", roomCode)
		return
	}

	switch msgType {
	case "join_game":
		m.handleJoinGame(game, sendToClientFunc)
	case "select_playlist":
		m.handleSelectPlaylist(game, msg, broadcastFunc)
	case "submit_answer":
		m.handleSubmitAnswer(game, username, msg, broadcastFunc)
	case "next_round":
		m.handleNextRound(game, broadcastFunc)
	}
}

func (m *BlindTestManager) handleSelectPlaylist(game *models.BlindTestGame, msg map[string]interface{}, broadcastFunc func(string, interface{})) {
	playlistID, ok := msg["playlist_id"].(string)
	if !ok {
		log.Printf("Playlist ID missing")
		return
	}

	if game.Status == "playing" && game.RoundNumber > 0 {
		log.Printf("Game already in progress for room %s", game.RoomCode)
		return
	}

	game.Config.PlaylistID = playlistID
	game.Status = "playing"

	m.StartGameRound(game, broadcastFunc)
}

func (m *BlindTestManager) StartGameRound(game *models.BlindTestGame, broadcastFunc func(string, interface{})) {
	game.RoundNumber++
	game.ResetRound()

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

	// Diffuse l'état de la manche à tous les clients WebSocket
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

	go func() {
		time.Sleep(time.Duration(game.Config.TimePerRound) * time.Second)
		m.endRound(game, broadcastFunc)
		m.checkGameEnd(game, broadcastFunc)
	}()
}

// Synchronise les réponses des joueurs avec la logique du serveur et recalcule les scores en temps réel
// Applique le multiplicateur de points (2x si titre+artiste, 1x sinon)
func (m *BlindTestManager) handleSubmitAnswer(game *models.BlindTestGame, username string, msg map[string]interface{}, broadcastFunc func(string, interface{})) {
	if game.CurrentRound == nil {
		return
	}

	if game.AnsweredPlayers[username] {
		return
	}

	trackAnswer, _ := msg["track"].(string)
	artistAnswer, _ := msg["artist"].(string)

	isCorrect := m.checkAnswer(trackAnswer, artistAnswer, game.CurrentRound)

	pointMultiplier := m.getPointMultiplier(trackAnswer, artistAnswer, game.CurrentRound)

	// Calculer le rang (position parmi les réponses correctes)
	rank := game.CurrentRound.CorrectAnswersCount + 1
	game.RecordAnswer(username, isCorrect, rank, pointMultiplier)

	log.Printf("Joueur %s a répondu: correct=%v, rang=%d, multiplicateur=%dx", username, isCorrect, rank, pointMultiplier)
}

func (m *BlindTestManager) checkAnswer(trackAnswer, artistAnswer string, round *models.BlindTestRound) bool {
	trackAnswer = strings.ToLower(strings.TrimSpace(trackAnswer))
	artistAnswer = strings.ToLower(strings.TrimSpace(artistAnswer))

	if trackAnswer == "" && artistAnswer == "" {
		return false
	}

	trackMatch := m.fuzzyMatch(trackAnswer, round.TrackName)
	artistMatch := m.fuzzyMatch(artistAnswer, round.ArtistName)

	return trackMatch || artistMatch
}

// Logique de scoring: 2x points si titre ET artiste trouvés, 1x sinon
func (m *BlindTestManager) getPointMultiplier(trackAnswer, artistAnswer string, round *models.BlindTestRound) int {
	trackAnswer = strings.ToLower(strings.TrimSpace(trackAnswer))
	artistAnswer = strings.ToLower(strings.TrimSpace(artistAnswer))

	trackMatch := m.fuzzyMatch(trackAnswer, round.TrackName)
	artistMatch := m.fuzzyMatch(artistAnswer, round.ArtistName)

	if trackMatch && artistMatch {
		return 2
	}
	return 1
}

// Fuzzy matching tolérant avec distance de Levenshtein (60% similitude minimum)
func (m *BlindTestManager) fuzzyMatch(answer, correct string) bool {
	answer = strings.ToLower(strings.TrimSpace(answer))
	correct = strings.ToLower(strings.TrimSpace(correct))

	if answer == "" {
		return false
	}

	if answer == correct {
		return true
	}

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

	if len(correct) >= 3 {
		similarity := levenshteinSimilarity(answer, correct)
		if similarity >= 0.6 {
			return true
		}
	}

	return false
}

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
				curr[j-1]+1,
				min(prev[j]+1,
					prev[j-1]+cost))
		}
		prev = curr
	}

	return prev[len(s2)]
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// Envoie le message round_end au WebSocket avec les résultats et scores mis à jour
func (m *BlindTestManager) endRound(game *models.BlindTestGame, broadcastFunc func(string, interface{})) {
	if game.CurrentRound == nil {
		return
	}

	roundScores := m.calculateRoundScores(game)

	game.RoundHistory = append(game.RoundHistory, *game.CurrentRound)

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

	broadcastFunc(game.RoomCode, map[string]interface{}{
		"type":   "scoreboard_update",
		"scores": game.GetScoreboard(),
	})
}

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

	m.StartGameRound(game, broadcastFunc)
}

func (m *BlindTestManager) RemoveGame(roomCode string) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	delete(m.games, roomCode)
	fmt.Printf("Blind test supprimé pour room %s\n", roomCode)
}

// Synchronise les nouveau clients avec l'état du jeu en cours (temps restant, round actuel)
func (m *BlindTestManager) handleJoinGame(game *models.BlindTestGame, sendToClientFunc func(interface{})) {
	if game == nil {
		log.Printf("Jeu non trouvé pour handleJoinGame")
		return
	}

	if game.Status == "playing" && game.CurrentRound != nil {
		elapsed := time.Since(game.CurrentRound.StartTime).Seconds()
		remaining := float64(game.Config.TimePerRound) - elapsed
		if remaining < 0 {
			remaining = 0
		}

		msg := map[string]interface{}{
			"type":           "round_start",
			"round_number":   game.RoundNumber,
			"total_rounds":   game.Config.NumRounds,
			"preview_url":    game.CurrentRound.PreviewURL,
			"cover_image":    game.CurrentRound.CoverImage,
			"duration":       game.Config.TimePerRound,
			"remaining_time": remaining,
			"track_name":     game.CurrentRound.TrackName,
			"artist_name":    game.CurrentRound.ArtistName,
		}
		sendToClientFunc(msg)
	} else if game.Status == "waiting_for_playlist" {
		// Le jeu attend qu'on sélectionne une playlist
		log.Printf("Jeu en attente de playlist pour room %s", game.RoomCode)
	}
}
