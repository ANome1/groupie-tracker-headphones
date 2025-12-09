package services

import (
	"database/sql"
	"log"
	"math/rand"
	"time"

	"groupie-tracker/database"
	"groupie-tracker/models"
)

type RoomService struct {
	DB *database.Database
}

func (rs *RoomService) GenerateRoomCode() string {
	const letters = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	for {
		code := make([]byte, 6)
		for i := range code {
			code[i] = letters[rand.Intn(len(letters))]
		}
		codeStr := string(code)

		// Vérifier que le code n'existe pas déjà
		var exists string
		err := rs.DB.DB.QueryRow("SELECT code FROM rooms WHERE code = ?", codeStr).Scan(&exists)
		if err == sql.ErrNoRows {

			return codeStr
		}
		if err != nil && err != sql.ErrNoRows {
			log.Printf("Error checking room code: %v", err)
			continue
		}
	}
}

func (rs *RoomService) CreateRoom(name, gameType string, hostID int) (*models.Room, error) {
	code := rs.GenerateRoomCode()

	result, err := rs.DB.DB.Exec(
		"INSERT INTO rooms (name, code, host_id, game_type, status) VALUES (?, ?, ?, ?, ?)",
		name, code, hostID, gameType, "waiting",
	)
	if err != nil {
		log.Printf("Error creating room: %v", err)
		return nil, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	return &models.Room{
		ID:        int(id),
		Name:      name,
		Code:      code,
		HostID:    hostID,
		GameType:  gameType,
		Status:    "waiting",
		CreatedAt: time.Now(),
	}, nil
}

func (rs *RoomService) GetRoomByCode(code string) (*models.Room, error) {
	var room models.Room
	err := rs.DB.DB.QueryRow(
		"SELECT id, name, code, host_id, game_type, max_players, status, created_at FROM rooms WHERE code = ?",
		code,
	).Scan(&room.ID, &room.Name, &room.Code, &room.HostID, &room.GameType, &room.MaxPlayers, &room.Status, &room.CreatedAt)

	if err != nil {
		return nil, err
	}
	return &room, nil
}

func (rs *RoomService) GetRoomByID(id int) (*models.Room, error) {
	var room models.Room
	err := rs.DB.DB.QueryRow(
		"SELECT id, name, code, host_id, game_type, max_players, status, created_at FROM rooms WHERE id = ?",
		id,
	).Scan(&room.ID, &room.Name, &room.Code, &room.HostID, &room.GameType, &room.MaxPlayers, &room.Status, &room.CreatedAt)

	if err != nil {
		return nil, err
	}
	return &room, nil
}

func (rs *RoomService) JoinRoom(roomID, userID int) error {
	_, err := rs.DB.DB.Exec(
		"INSERT INTO room_participants (room_id, user_id) VALUES (?, ?)",
		roomID, userID,
	)
	if err != nil {
		log.Printf("Error joining room: %v", err)
		return err
	}
	return nil
}

func (rs *RoomService) GetRoomParticipants(roomID int) ([]models.RoomParticipant, error) {
	var participants []models.RoomParticipant

	rows, err := rs.DB.DB.Query(
		"SELECT id, room_id, user_id, score, joined_at FROM room_participants WHERE room_id = ?",
		roomID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var participant models.RoomParticipant
		err := rows.Scan(&participant.ID, &participant.RoomID, &participant.UserID, &participant.Score, &participant.JoinedAt)
		if err != nil {
			return nil, err
		}
		participants = append(participants, participant)
	}

	return participants, nil
}

func (rs *RoomService) GetRoomParticipantsWithUsers(roomID int) ([]models.RoomParticipantWithUser, error) {
	var participants []models.RoomParticipantWithUser

	rows, err := rs.DB.DB.Query(
		`SELECT rp.id, rp.room_id, rp.user_id, u.username, rp.score, rp.joined_at 
		FROM room_participants rp
		JOIN users u ON rp.user_id = u.id
		WHERE rp.room_id = ?`,
		roomID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var participant models.RoomParticipantWithUser
		err := rows.Scan(&participant.ID, &participant.RoomID, &participant.UserID, &participant.Username, &participant.Score, &participant.JoinedAt)
		if err != nil {
			return nil, err
		}
		participants = append(participants, participant)
	}

	return participants, nil
}

func (rs *RoomService) LeaveRoom(roomID, userID int) error {
	_, err := rs.DB.DB.Exec(
		"DELETE FROM room_participants WHERE room_id = ? AND user_id = ?",
		roomID, userID,
	)
	if err != nil {
		log.Printf("Error leaving room: %v", err)
		return err
	}
	return nil
}

func (rs *RoomService) DeleteRoom(roomID int) error {
	_, err := rs.DB.DB.Exec("DELETE FROM rooms WHERE id = ?", roomID)
	if err != nil {
		log.Printf("Error deleting room: %v", err)
		return err
	}
	return nil
}
