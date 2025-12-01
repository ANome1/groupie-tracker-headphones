package models

// RESPONSABLE: @Nome (base), @Quoc Huy (Blind Test), @ilian (Petit Bac)
// Modèle de salle de jeu

import "time"

type Room struct {
	ID         int       `json:"id"`
	Name       string    `json:"name"`
	Code       string    `json:"code"` // Code à 6 caractères (ex: ABC123)
	HostID     int       `json:"host_id"`
	GameType   string    `json:"game_type"` // "blindtest" ou "petitbac"
	MaxPlayers int       `json:"max_players"`
	Status     string    `json:"status"` // "waiting", "in_progress", "finished"
	CreatedAt  time.Time `json:"created_at"`
}

type RoomParticipant struct {
	ID       int       `json:"id"`
	RoomID   int       `json:"room_id"`
	UserID   int       `json:"user_id"`
	Username string    `json:"username"` // Pour affichage
	Score    int       `json:"score"`
	JoinedAt time.Time `json:"joined_at"`
}

// TODO @Nome: Ajouter les méthodes suivantes:
// - func CreateRoom(name, gameType string, hostID int) (*Room, error)
// - func GetRoomByCode(code string) (*Room, error)
// - func JoinRoom(roomID, userID int) error
// - func GetRoomParticipants(roomID int) ([]*RoomParticipant, error)
// - func GenerateRoomCode() string // Génère un code aléatoire de 6 caractères
