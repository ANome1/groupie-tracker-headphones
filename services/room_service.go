package services

// TODO @Nome: CreateRoom, GetRoomByCode, JoinRoom, LeaveRoom, GenerateRoomCode
import (
	"groupie-tracker/database"
	"math/rand"
)

type RoomService struct {
	DB *database.Database
}

func GenerateRoomCode() string {
	const letters = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	code := make([]byte, 6)
	for i := range code {
		code[i] = letters[rand.Intn(len(letters))]
	}
	return string(code)
}
