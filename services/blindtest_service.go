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
	track, err := GetRandomTrackFromDeezerPlaylist(game.Config.PlaylistID)
	if err != nil {
		log.Printf("Erreur récupération track: %v", err)
		broadcastFunc(game.RoomCode, map[string]interface{}{
			"type":  "error",
			"error": "Impossible de charger la musique",
		})
		return
	}

	// Créer la manche
	game.CurrentRound = &models.BlindTestRound{
		RoundNumber: game.RoundNumber,
		TrackID:     track.ID,
		TrackName:   track.Title,
		ArtistName:  track.Artist.Name,
		PreviewURL:  track.Preview,
		CoverImage:  track.Album.CoverMedium,
		Duration:    game.Config.TimePerRound,
		StartTime:   time.Now(),
	}

	// Broadcast aux clients
	broadcastFunc(game.RoomCode, map[string]interface{}{
		"type":         "round_start",
		"round_number": game.RoundNumber,
		"total_rounds": game.Config.NumRounds,
		"preview_url":  track.Preview,
		"duration":     game.Config.TimePerRound,
	})

	// Timer automatique pour finir la manche
	go func() {
		time.Sleep(time.Duration(game.Config.TimePerRound) * time.Second)
		m.endRound(game, broadcastFunc)
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

	// Calcul du temps écoulé
	elapsed := time.Since(game.CurrentRound.StartTime).Seconds()

	// Vérifier si la réponse est correcte
	isCorrect := m.checkAnswer(trackAnswer, artistAnswer, game.CurrentRound)

	// Calcul du bonus de temps (max 50 points)
	timeBonus := 0
	if isCorrect {
		remainingTime := float64(game.Config.TimePerRound) - elapsed
		timeBonus = int(remainingTime * 50 / float64(game.Config.TimePerRound))
	}

	// Enregistrer la réponse
	game.RecordAnswer(username, isCorrect, timeBonus)

	log.Printf("Joueur %s a répondu: correct=%v, bonus=%d", username, isCorrect, timeBonus)
}

// checkAnswer vérifie si la réponse est correcte
func (m *BlindTestManager) checkAnswer(trackAnswer, artistAnswer string, round *models.BlindTestRound) bool {
	trackMatch := m.fuzzyMatch(trackAnswer, round.TrackName)
	artistMatch := m.fuzzyMatch(artistAnswer, round.ArtistName)

	// Au moins l'un des deux doit être correct
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

	// Match partiel (au moins 70% de similarité)
	return strings.Contains(correct, answer) || strings.Contains(answer, correct)
}

// endRound termine la manche et envoie les résultats
func (m *BlindTestManager) endRound(game *models.BlindTestGame, broadcastFunc func(string, interface{})) {
	if game.CurrentRound == nil {
		return
	}

	// Ajouter à l'historique
	game.RoundHistory = append(game.RoundHistory, *game.CurrentRound)

	// Broadcast les résultats
	broadcastFunc(game.RoomCode, map[string]interface{}{
		"type":           "round_end",
		"correct_track":  game.CurrentRound.TrackName,
		"correct_artist": game.CurrentRound.ArtistName,
		"cover_image":    game.CurrentRound.CoverImage,
		"scores":         game.GetScoreboard(),
	})

	// Mettre à jour le scoreboard
	broadcastFunc(game.RoomCode, map[string]interface{}{
		"type":   "scoreboard_update",
		"scores": game.GetScoreboard(),
	})

	// Vérifier si c'est la fin du jeu
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
