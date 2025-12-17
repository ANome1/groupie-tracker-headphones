package models

import "time"

type GameSession struct {
	ID        int
	GameType  string
	Players   []int
	RoomID    int
	GameData  string
	StartedAt time.Time
	EndedAt   *time.Time
}
