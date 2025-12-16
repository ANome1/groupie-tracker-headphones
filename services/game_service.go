package services

// TODO @Quoc Huy & @ilian: StartGame, EndGame, GetGameState, SaveGameState

import (
	"groupie-tracker/models"
	"log"
	"math/rand"
	"sync"
	"time"
) // PetitBacManager manages multiple game sessions
type PetitBacManager struct {
	Games map[string]*models.PetitBacGame
	mutex sync.RWMutex
}

// Global manager instance
var Manager = &PetitBacManager{
	Games: make(map[string]*models.PetitBacGame),
}

// Global BlindTest manager instance
var BlindTestMgr = NewBlindTestManager()

// CreateGame initializes a new game session
func (m *PetitBacManager) CreateGame(roomID, hostID string, players []string, playerNames map[string]string) *models.PetitBacGame {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	gameID := roomID // Using roomID as gameID for simplicity
	config := models.PetitBacConfig{
		Categories:   []string{"Artiste", "Groupe de musique", "Album", "Instrument", "Featuring"},
		TimePerRound: 60,
		NumRounds:    9,
	}

	game := models.NewPetitBacGame(gameID, roomID, hostID, players, playerNames, config)
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

	game.Lock()
	defer game.Unlock()

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

	game.Lock()
	defer game.Unlock()

	currentRound := game.Rounds[len(game.Rounds)-1]
	currentRound.Responses[playerID] = &models.PlayerResponse{
		PlayerID:   playerID,
		Answers:    answers,
		Submitted:  true,
		SubmitTime: time.Now(),
	}

	// Check if all players submitted to trigger validation phase early
	log.Printf("SubmitAnswers: Player %s submitted. Total responses: %d/%d", playerID, len(currentRound.Responses), len(game.Players))
	if len(currentRound.Responses) == len(game.Players) {
		game.State = models.GameStateVoting
		// Trigger validation phase (caller handles broadcast)
	}
}

// SubmitVote records a vote
func (m *PetitBacManager) SubmitVote(gameID string, vote models.Vote) {
	game := m.GetGame(gameID)
	if game == nil || len(game.Rounds) == 0 {
		return
	}

	game.Lock()
	defer game.Unlock()

	currentRound := game.Rounds[len(game.Rounds)-1]
	currentRound.Votes = append(currentRound.Votes, vote)

	// TODO: Check if voting is complete (all players voted for all answers)
	// For simplicity, we might just rely on a timer or manual "Next Round"
	// But let's calculate scores on the fly or when requested
}

// CalculateScores updates the scores based on votes and answers
func (m *PetitBacManager) CalculateScores(gameID string) map[string]int {
	game := m.GetGame(gameID)
	if game == nil || len(game.Rounds) == 0 {
		return nil
	}

	game.RLock()
	defer game.RUnlock()

	currentRound := game.Rounds[len(game.Rounds)-1]

	// Initialize validity map (default false - requires validation)
	validity := make(map[string]map[string]bool)
	for _, pID := range game.Players {
		validity[pID] = make(map[string]bool)
		if resp, ok := currentRound.Responses[pID]; ok {
			for cat := range resp.Answers {
				validity[pID][cat] = false
			}
		}
	}

	// Count votes
	type voteCount struct {
		valid   int
		invalid int
	}
	votes := make(map[string]map[string]*voteCount)

	for _, v := range currentRound.Votes {
		if _, ok := votes[v.TargetPlayer]; !ok {
			votes[v.TargetPlayer] = make(map[string]*voteCount)
		}
		if _, ok := votes[v.TargetPlayer][v.Category]; !ok {
			votes[v.TargetPlayer][v.Category] = &voteCount{}
		}

		if v.IsValid {
			votes[v.TargetPlayer][v.Category].valid++
		} else {
			votes[v.TargetPlayer][v.Category].invalid++
		}
	}

	// Apply majority rule: if valid > invalid, mark as valid
	for pID, cats := range votes {
		for cat, count := range cats {
			if count.valid > count.invalid {
				if validity[pID] != nil {
					validity[pID][cat] = true
				}
			}
		}
	}

	// Check for uniqueness
	answerCounts := make(map[string]map[string]int)

	for _, pID := range game.Players {
		if resp, ok := currentRound.Responses[pID]; ok {
			for cat, ans := range resp.Answers {
				if ans == "" {
					continue
				}
				if _, ok := answerCounts[cat]; !ok {
					answerCounts[cat] = make(map[string]int)
				}
				answerCounts[cat][ans]++
			}
		}
	}

	projectedScores := make(map[string]int)
	for k, v := range game.Scores {
		projectedScores[k] = v
	}

	for _, pID := range game.Players {
		points := 0
		if resp, ok := currentRound.Responses[pID]; ok {
			for cat, ans := range resp.Answers {
				// Check if category is valid
				isValidCategory := false
				for _, c := range game.Config.Categories {
					if c == cat {
						isValidCategory = true
						break
					}
				}
				if !isValidCategory {
					continue
				}

				if ans == "" {
					continue
				}
				if !validity[pID][cat] {
					continue
				}

				// Points accumulation: 2 points for validated answer
				points += 2
			}
		}
		projectedScores[pID] += points
	}

	return projectedScores
}

// NextRound finalizes the current round and starts the next one
func (m *PetitBacManager) NextRound(gameID string) (*models.RoundUpdate, bool) {
	// Calculate final scores for the current round
	scores := m.CalculateScores(gameID)

	game := m.GetGame(gameID)
	if game == nil {
		return nil, false
	}

	game.Lock()
	defer game.Unlock()

	// Update scores
	if scores != nil {
		game.Scores = scores
	}

	// Check if game is finished
	if game.CurrentRound >= game.Config.NumRounds {
		game.State = models.GameStateFinished
		return nil, true // Game Over
	}

	// Start new round logic
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
	}, false
}
