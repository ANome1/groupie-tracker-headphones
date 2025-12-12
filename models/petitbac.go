package models

import (
	"encoding/json"
	"time"
)

// TODO @ilian: Constante NbrsManche = 9
// TODO @ilian: Structures PetitBacGame, PetitBacAnswer
// TODO @ilian: Variable scoreboardActualPointInGame
// TODO @ilian: Méthodes de l'interface Game

// Game States
const (
	GameStateWaiting  = "WAITING"
	GameStatePlaying  = "PLAYING"
	GameStateVoting   = "VOTING"
	GameStateFinished = "FINISHED"
)

// Defaults
const (
	DefaultRounds    = 5
	DefaultRoundTime = 60 // seconds
)

// PetitBacGame represents the complete state of a Petit Bac game session.
type PetitBacGame struct {
	ID           string
	RoomID       string
	State        string
	Config       PetitBacConfig
	CurrentRound int
	UsedLetters  []rune
	Rounds       []*PetitBacRound
	Scores       map[string]int // PlayerID -> Score
	Players      []string       // List of PlayerIDs
	CreatedAt    time.Time
}

// PetitBacConfig defines the rules for the game.
type PetitBacConfig struct {
	Categories   []string
	TimePerRound int
	NumRounds    int
}

// PetitBacRound stores data for a specific round.
type PetitBacRound struct {
	RoundNumber int
	Letter      rune
	StartTime   time.Time
	EndTime     time.Time
	Responses   map[string]*PlayerResponse // PlayerID -> Response
	Votes       []Vote                     // List of all votes in this round
}

// PlayerResponse represents a player's answers for a round.
type PlayerResponse struct {
	PlayerID   string
	Answers    map[string]string // Category -> Answer
	Submitted  bool
	SubmitTime time.Time
}

// Vote represents a validation vote from one player for another player's answer.
type Vote struct {
	VoterID      string
	TargetPlayer string
	Category     string
	IsValid      bool
}

// MessageOut is a helper for WebSocket responses.
type MessageOut struct {
	Type string      // Ex: "GAME_START", "NEW_ROUND", "TIME_UPDATE", "SCORES"
	Data interface{} // Contient la structure Go appropriée
}

// RoundUpdate is a helper for sending round info.
type RoundUpdate struct {
	Letter      string
	Duration    int
	RoundNumber int
}

// MessageIn represents the standard format for messages received from the client.
type MessageIn struct {
	Type string          // Ex: "SUBMIT_ANSWERS", "SUBMIT_VOTE"
	Data json.RawMessage // Le contenu brut qui sera décodé selon le Type
}

// AnswersPayload is the data expected when Type is "SUBMIT_ANSWERS".
type AnswersPayload struct {
	PlayerID string
	Answers  map[string]string
}

// NewPetitBacGame creates a new game instance with default or provided config.
func NewPetitBacGame(id, roomID string, players []string, config PetitBacConfig) *PetitBacGame {
	if config.NumRounds == 0 {
		config.NumRounds = DefaultRounds
	}
	if config.TimePerRound == 0 {
		config.TimePerRound = DefaultRoundTime
	}
	if len(config.Categories) == 0 {
		// Default categories if none provided
		config.Categories = []string{"Pays", "Ville", "Animal", "Métier", "Objet", "Prénom"}
	}

	// Initialize scores
	scores := make(map[string]int)
	for _, p := range players {
		scores[p] = 0
	}

	return &PetitBacGame{
		ID:           id,
		RoomID:       roomID,
		State:        GameStateWaiting,
		Config:       config,
		CurrentRound: 0,
		UsedLetters:  make([]rune, 0),
		Rounds:       make([]*PetitBacRound, 0),
		Scores:       scores,
		Players:      players,
		CreatedAt:    time.Now(),
	}
}
