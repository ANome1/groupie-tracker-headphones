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

// BlindTestGame représente une session de jeu
type BlindTestGame struct {
	RoomCode string
	HostID   int
	Status   string // "waiting", "playing", "round_end", "finished"

	Config       BlindTestConfig
	CurrentRound *BlindTestRound
	RoundHistory []BlindTestRound
	RoundNumber  int

	Scores          map[string]int    // username -> score
	AnsweredPlayers map[string]bool   // username -> has answered this round
	PlayerNames     map[string]string // playerID -> username

	Timer *time.Timer // Timer pour la fin de la manche

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
		PlayerNames:     make(map[string]string),
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

// RecordAnswer enregistre qu'un joueur a répondu correctement
// pointMultiplier permet de multiplier les points (ex: 2 si titre + artiste trouvés)
func (g *BlindTestGame) RecordAnswer(username string, isCorrect bool, rank int, pointMultiplier int) {
	g.mutex.Lock()
	defer g.mutex.Unlock()

	if g.AnsweredPlayers[username] {
		return // Already answered
	}

	g.AnsweredPlayers[username] = true

	if isCorrect && g.CurrentRound != nil {
		// Points décroissants selon la position: 1er=100, 2e=80, 3e=60, 4e=40, 5e=20
		pointsByRank := []int{100, 80, 60, 40, 20}
		points := 0
		if rank-1 < len(pointsByRank) {
			points = pointsByRank[rank-1]
		}
		// Appliquer le multiplicateur (ex: 2x si titre + artiste trouvés)
		if pointMultiplier < 1 {
			pointMultiplier = 1
		}
		points *= pointMultiplier

		g.Scores[username] += points
		g.CurrentRound.CorrectAnswersCount++
		g.CurrentRound.CorrectAnswersPlayers = append(g.CurrentRound.CorrectAnswersPlayers, username)
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
