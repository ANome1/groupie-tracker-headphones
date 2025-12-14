package services

// TODO @Quoc Huy & @ilian: StartGame, EndGame, GetGameState, SaveGameState

import (
	"groupie-tracker/models"
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

// SubmitVote records a vote
func (m *PetitBacManager) SubmitVote(gameID string, vote models.Vote) {
	game := m.GetGame(gameID)
	if game == nil || len(game.Rounds) == 0 {
		return
	}

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

	currentRound := game.Rounds[len(game.Rounds)-1]

	// Initialize validity map (default true)
	validity := make(map[string]map[string]bool)
	for _, pID := range game.Players {
		validity[pID] = make(map[string]bool)
		if resp, ok := currentRound.Responses[pID]; ok {
			for cat := range resp.Answers {
				validity[pID][cat] = true
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

	// Apply majority rule: if invalid > valid, mark as invalid
	for pID, cats := range votes {
		for cat, count := range cats {
			if count.invalid > count.valid {
				if validity[pID] != nil {
					validity[pID][cat] = false
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
				if ans == "" {
					continue
				}
				if !validity[pID][cat] {
					continue
				}

				count := answerCounts[cat][ans]
				if count == 1 {
					points += 2
				} else {
					points += 1
				}
			}
		}
		projectedScores[pID] += points
	}

	return projectedScores
}
