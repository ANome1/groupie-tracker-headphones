package models

import (
	"sync"
	"time"
)

type BlindTestConfig struct {
	PlaylistID   string
	TimePerRound int
	NumRounds    int
}

type BlindTestRound struct {
	RoundNumber           int
	TrackID               int
	TrackName             string
	ArtistName            string
	PreviewURL            string
	CoverImage            string
	Duration              int
	StartTime             time.Time
	CorrectAnswersCount   int
	CorrectAnswersPlayers []string
}

type BlindTestGame struct {
	RoomCode string
	HostID   int
	Status   string

	Config       BlindTestConfig
	CurrentRound *BlindTestRound
	RoundHistory []BlindTestRound
	RoundNumber  int

	Scores          map[string]int
	AnsweredPlayers map[string]bool
	PlayerNames     map[string]string
	Timer           *time.Timer

	mutex sync.RWMutex
}

func NewBlindTestGame(roomCode string, hostID int, config BlindTestConfig) *BlindTestGame {
	return &BlindTestGame{
		RoomCode:        roomCode,
		HostID:          hostID,
		Status:          "waiting",
		Config:          config,
		Scores:          make(map[string]int),
		AnsweredPlayers: make(map[string]bool),
		PlayerNames:     make(map[string]string),
		RoundNumber:     0,
		RoundHistory:    []BlindTestRound{},
	}
}

// AddPlayer registers a player in the game
func (g *BlindTestGame) AddPlayer(username string) {
	g.mutex.Lock()
	defer g.mutex.Unlock()

	if _, exists := g.Scores[username]; !exists {
		g.Scores[username] = 0
	}
}

// Enregistre la réponse avec le rang et le multiplicateur (2x si titre+artiste)
func (g *BlindTestGame) RecordAnswer(username string, isCorrect bool, rank int, pointMultiplier int) {
	g.mutex.Lock()
	defer g.mutex.Unlock()

	if g.AnsweredPlayers[username] {
		return
	}

	g.AnsweredPlayers[username] = true

	if isCorrect && g.CurrentRound != nil {
		pointsByRank := []int{100, 80, 60, 40, 20}
		points := 0
		if rank-1 < len(pointsByRank) {
			points = pointsByRank[rank-1]
		}
		if pointMultiplier < 1 {
			pointMultiplier = 1
		}
		points *= pointMultiplier

		g.Scores[username] += points
		g.CurrentRound.CorrectAnswersCount++
		g.CurrentRound.CorrectAnswersPlayers = append(g.CurrentRound.CorrectAnswersPlayers, username)
	}
}

// ResetRound prepares for the next round
func (g *BlindTestGame) ResetRound() {
	g.mutex.Lock()
	defer g.mutex.Unlock()

	g.AnsweredPlayers = make(map[string]bool)
}

// GetScoreboard returns current scores
func (g *BlindTestGame) GetScoreboard() map[string]int {
	g.mutex.RLock()
	defer g.mutex.RUnlock()

	scoreboard := make(map[string]int)
	for username, score := range g.Scores {
		scoreboard[username] = score
	}
	return scoreboard
}
