package models

// TODO @Nome: Structure Room (ID, Name, Code, HostID, GameType, MaxPlayers, Status, CreatedAt)
// TODO @Nome: Structure RoomParticipant
// TODO @Nome: Fonctions CRUD
import "time"

type Room struct {
	ID         int
	Name       string
	Code       string
	HostID     int
	GameType   string
	MaxPlayers int
	Status     string
	CreatedAt  time.Time
}

type RoomParticipant struct {
	ID       int
	RoomID   int
	UserID   int
	Score    int
	JoinedAt time.Time
}
