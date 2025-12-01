package services

// RESPONSABLE: @Nome
// Service pour la gestion des salles de jeu

import "math/rand"

// RoomService - Gestion des salles
type RoomService struct {
	// TODO @Nome: Ajouter la connexion à la base de données
}

// NewRoomService - Crée une nouvelle instance du service
func NewRoomService() *RoomService {
	return &RoomService{}
}

// CreateRoom - Crée une nouvelle salle
// TODO @Nome:
// - Générer un code unique (GenerateRoomCode)
// - Insérer la salle en base de données
// - Ajouter l'hôte comme participant
// - Retourner la salle créée
func (s *RoomService) CreateRoom(name, gameType string, hostID, maxPlayers int) (*Room, error) {
	// TODO: Implementation
	return nil, nil
}

// GenerateRoomCode - Génère un code aléatoire de 6 caractères
// TODO @Nome:
// - Format: ABCDEF (lettres majuscules et chiffres)
// - Vérifier l'unicité en base de données
func (s *RoomService) GenerateRoomCode() string {
	const charset = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	code := make([]byte, 6)
	for i := range code {
		code[i] = charset[rand.Intn(len(charset))]
	}
	return string(code)
}

// GetRoomByCode - Récupère une salle par son code
// TODO @Nome:
func (s *RoomService) GetRoomByCode(code string) (*Room, error) {
	// TODO: Implementation
	return nil, nil
}

// JoinRoom - Ajoute un joueur à une salle
// TODO @Nome:
// - Vérifier que la salle existe
// - Vérifier que la salle n'est pas pleine
// - Vérifier que le joueur n'est pas déjà dans la salle
// - Ajouter le participant
func (s *RoomService) JoinRoom(roomID, userID int) error {
	// TODO: Implementation
	return nil
}

// GetRoomParticipants - Liste tous les participants d'une salle
// TODO @Nome:
func (s *RoomService) GetRoomParticipants(roomID int) ([]*Participant, error) {
	// TODO: Implementation
	return nil, nil
}

// UpdateRoomStatus - Change le statut d'une salle
// TODO @Nome:
// - Statuts possibles: "waiting", "in_progress", "finished"
func (s *RoomService) UpdateRoomStatus(roomID int, status string) error {
	// TODO: Implementation
	return nil
}

type Room struct {
	ID         int
	Name       string
	Code       string
	HostID     int
	GameType   string
	MaxPlayers int
	Status     string
}

type Participant struct {
	ID       int
	RoomID   int
	UserID   int
	Username string
	Score    int
}
