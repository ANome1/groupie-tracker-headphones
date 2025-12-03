package models

// TODO @Quoc Huy & @ilian: Structure GameSession
// TODO @Quoc Huy & @ilian: Interface Game
import "time"

type GameSession struct {
	ID        int
	RoomID    int
	GameData  string
	StartedAt time.Time
	EndedAt   *time.Time
}
