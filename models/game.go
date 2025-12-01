package models

// RESPONSABLE: @Quoc Huy & @ilian
// Interface commune pour les jeux

import "time"

type GameSession struct {
	ID        int        `json:"id"`
	RoomID    int        `json:"room_id"`
	GameData  string     `json:"game_data"` // JSON de l'état du jeu
	StartedAt time.Time  `json:"started_at"`
	EndedAt   *time.Time `json:"ended_at,omitempty"`
}

// Interface que doivent implémenter BlindTestGame et PetitBacGame
type Game interface {
	GetType() string        // Retourne "blindtest" ou "petitbac"
	IsFinished() bool       // Vérifie si le jeu est terminé
	GetScores() map[int]int // Retourne les scores (userID -> score)
}
