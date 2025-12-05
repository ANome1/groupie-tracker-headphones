package models

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
