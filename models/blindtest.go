package models

import (
	"sync"
	"time"
)

// BlindTestConfig configuration pour le blind test
type BlindTestConfig struct {
	PlaylistID   string
	TimePerRound int
	NumRounds    int
}

// BlindTestRound représente une manche en cours
type BlindTestRound struct {
	RoundNumber int
	TrackID     int
	TrackName   string
	ArtistName  string
	PreviewURL  string
	CoverImage  string
	Duration    int
	StartTime   time.Time
}

// BlindTestGame représente une session de jeu
type BlindTestGame struct {
	RoomCode string
	HostID   int
	Status   string // "waiting", "playing", "round_end", "finished"

	Config       BlindTestConfig
	CurrentRound *BlindTestRound
	RoundHistory []BlindTestRound
	RoundNumber  int

	Scores          map[string]int  // username -> score
	AnsweredPlayers map[string]bool // username -> has answered this round

	mutex sync.RWMutex
}

// NewBlindTestGame crée une nouvelle partie
func NewBlindTestGame(roomCode string, hostID int, config BlindTestConfig) *BlindTestGame {
	return &BlindTestGame{
		RoomCode:        roomCode,
		HostID:          hostID,
		Status:          "waiting",
		Config:          config,
		Scores:          make(map[string]int),
		AnsweredPlayers: make(map[string]bool),
		RoundNumber:     0,
		RoundHistory:    []BlindTestRound{},
	}
}

// AddPlayer ajoute un joueur au jeu
func (g *BlindTestGame) AddPlayer(username string) {
	g.mutex.Lock()
	defer g.mutex.Unlock()

	if _, exists := g.Scores[username]; !exists {
		g.Scores[username] = 0
	}
}

// RecordAnswer enregistre qu'un joueur a répondu
func (g *BlindTestGame) RecordAnswer(username string, isCorrect bool, timeBonus int) {
	g.mutex.Lock()
	defer g.mutex.Unlock()

	if g.AnsweredPlayers[username] {
		return // Already answered
	}

	g.AnsweredPlayers[username] = true

	if isCorrect {
		// Base score + time bonus
		g.Scores[username] += 100 + timeBonus
	}
}

// ResetRound prépare pour la prochaine manche
func (g *BlindTestGame) ResetRound() {
	g.mutex.Lock()
	defer g.mutex.Unlock()

	g.AnsweredPlayers = make(map[string]bool)
}

// GetScoreboard retourne le classement actuel
func (g *BlindTestGame) GetScoreboard() map[string]int {
	g.mutex.RLock()
	defer g.mutex.RUnlock()

	scoreboard := make(map[string]int)
	for username, score := range g.Scores {
		scoreboard[username] = score
	}
	return scoreboard
}
