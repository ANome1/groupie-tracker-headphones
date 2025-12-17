package models

import (
	"encoding/json"
	"sync"
	"time"
)

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
	sync.RWMutex // Embed mutex for thread safety

	ID           string
	RoomID       string
	HostID       string // ID of the host player
	State        string
	Config       PetitBacConfig
	CurrentRound int
	UsedLetters  []rune
	Rounds       []*PetitBacRound
	Scores       map[string]int    // PlayerID -> Score
	Players      []string          // List of PlayerIDs
	PlayerNames  map[string]string // PlayerID -> Username
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
	RoundNumber int                        `json:"RoundNumber"`
	Letter      rune                       `json:"Letter"`
	StartTime   time.Time                  `json:"StartTime"`
	EndTime     time.Time                  `json:"EndTime"`
	Responses   map[string]*PlayerResponse `json:"Responses"`
	Votes       []Vote                     `json:"Votes"`
}

// PlayerResponse represents a player's answers for a round.
type PlayerResponse struct {
	PlayerID   string            `json:"PlayerID"`
	Answers    map[string]string `json:"Answers"`
	Submitted  bool              `json:"Submitted"`
	SubmitTime time.Time         `json:"SubmitTime"`
}

// Vote represents a validation vote from one player for another player's answer.
type Vote struct {
	VoterID      string `json:"VoterID"`
	TargetPlayer string `json:"TargetPlayer"`
	Category     string `json:"Category"`
	IsValid      bool   `json:"IsValid"`
}

// MessageOut is a helper for WebSocket responses.
type MessageOut struct {
	Type string      // Ex: "GAME_START", "NEW_ROUND", "TIME_UPDATE", "SCORES"
	Data interface{} // Contient la structure Go appropriée
}

// RoundUpdate is a helper for sending round info.
type RoundUpdate struct {
	Letter      string `json:"Letter"`
	Duration    int    `json:"Duration"`
	RoundNumber int    `json:"RoundNumber"`
}

// RoundResults stores the scoring results for a round
type RoundResults struct {
	RoundNumber int                                 `json:"RoundNumber"`
	Answers     map[string]map[string]*AnswerResult `json:"Answers"`     // PlayerID -> CategoryID -> Answer details
	RoundScores map[string]int                      `json:"RoundScores"` // PlayerID -> Points gained this round
	TotalScores map[string]int                      `json:"TotalScores"` // PlayerID -> Cumulative total
}

// AnswerResult stores details about an answer
type AnswerResult struct {
	Answer string `json:"Answer"`
	Points int    `json:"Points"` // 0, 1, or 2
	Valid  bool   `json:"Valid"`  // Was it validated?
	Unique bool   `json:"Unique"` // Is it unique?
}

// MessageIn represents the standard format for messages received from the client.
type MessageIn struct {
	Type string          `json:"Type"`
	Data json.RawMessage `json:"Data"`
}

// AnswersPayload is the data expected when Type is "SUBMIT_ANSWERS".
type AnswersPayload struct {
	PlayerID string            `json:"PlayerID"`
	Answers  map[string]string `json:"Answers"`
}

// NewPetitBacGame creates a new game instance with default or provided config.
func NewPetitBacGame(id, roomID, hostID string, players []string, playerNames map[string]string, config PetitBacConfig) *PetitBacGame {
	if config.NumRounds == 0 {
		config.NumRounds = DefaultRounds
	}
	if config.TimePerRound == 0 {
		config.TimePerRound = DefaultRoundTime
	}
	if len(config.Categories) == 0 {
		// Default categories if none provided
		config.Categories = []string{"Artiste", "Groupe de musique", "Album", "Instrument", "Featuring"}
	}

	// Initialize scores
	scores := make(map[string]int)
	for _, p := range players {
		scores[p] = 0
	}

	return &PetitBacGame{
		ID:           id,
		RoomID:       roomID,
		HostID:       hostID,
		State:        GameStateWaiting,
		Config:       config,
		CurrentRound: 0,
		UsedLetters:  make([]rune, 0),
		Rounds:       make([]*PetitBacRound, 0),
		Scores:       scores,
		Players:      players,
		PlayerNames:  playerNames,
		CreatedAt:    time.Now(),
	}
}
