package services

import (
	"groupie-tracker/models"
	"math/rand"
	"sync"
	"time"
)

// PetitBacManager manages multiple game sessions
type PetitBacManager struct {
	Games map[string]*models.PetitBacGame
	mutex sync.RWMutex
}

// Global manager instance
var Manager = &PetitBacManager{
	Games: make(map[string]*models.PetitBacGame),
}

// CreateGame initializes a new game session
func (m *PetitBacManager) CreateGame(roomID string, players []string) *models.PetitBacGame {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	gameID := roomID // Using roomID as gameID for simplicity
	config := models.PetitBacConfig{
		Categories:   []string{"Pays", "Ville", "Animal", "Métier", "Objet", "Prénom"},
		TimePerRound: 60,
		NumRounds:    9,
	}

	game := models.NewPetitBacGame(gameID, roomID, players, config)
	m.Games[gameID] = game
	return game
}

// GetGame retrieves a game by ID
func (m *PetitBacManager) GetGame(gameID string) *models.PetitBacGame {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	return m.Games[gameID]
}

// StartRound initiates a new round
func (m *PetitBacManager) StartRound(gameID string) *models.RoundUpdate {
	game := m.GetGame(gameID)
	if game == nil {
		return nil
	}

	// Pick a random letter not used yet
	letters := "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	var letter rune
	for {
		letter = rune(letters[rand.Intn(len(letters))])
		used := false
		for _, l := range game.UsedLetters {
			if l == letter {
				used = true
				break
			}
		}
		if !used {
			break
		}
		// Safety break if all letters used (unlikely with 9 rounds)
		if len(game.UsedLetters) >= 26 {
			break
		}
	}

	game.UsedLetters = append(game.UsedLetters, letter)
	game.CurrentRound++

	round := &models.PetitBacRound{
		RoundNumber: game.CurrentRound,
		Letter:      letter,
		StartTime:   time.Now(),
		EndTime:     time.Now().Add(time.Duration(game.Config.TimePerRound) * time.Second),
		Responses:   make(map[string]*models.PlayerResponse),
	}

	game.Rounds = append(game.Rounds, round)
	game.State = models.GameStatePlaying

	return &models.RoundUpdate{
		Letter:      string(letter),
		Duration:    game.Config.TimePerRound,
		RoundNumber: game.CurrentRound,
	}
}

// SubmitAnswers records a player's answers
func (m *PetitBacManager) SubmitAnswers(gameID, playerID string, answers map[string]string) {
	game := m.GetGame(gameID)
	if game == nil || len(game.Rounds) == 0 {
		return
	}

	currentRound := game.Rounds[len(game.Rounds)-1]
	currentRound.Responses[playerID] = &models.PlayerResponse{
		PlayerID:   playerID,
		Answers:    answers,
		Submitted:  true,
		SubmitTime: time.Now(),
	}

	// Check if all players submitted to trigger validation phase early
	if len(currentRound.Responses) == len(game.Players) {
		game.State = models.GameStateVoting
		// Trigger validation phase (caller handles broadcast)
	}
}
