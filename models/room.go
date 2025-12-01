package models

// RESPONSABLE: @Nome (base), @Quoc Huy (Blind Test), @ilian (Petit Bac)
// Modèle de salle de jeu

// TODO @Nome: Créer la structure Room avec:
// - ID int
// - Name string
// - Code string (6 caractères, ex: ABC123)
// - HostID int
// - GameType string ("blindtest" ou "petitbac")
// - MaxPlayers int
// - Status string ("waiting", "in_progress", "finished")
// - CreatedAt time.Time

// TODO @Nome: Créer la structure RoomParticipant avec:
// - ID int
// - RoomID int
// - UserID int
// - Username string
// - Score int
// - JoinedAt time.Time

// TODO @Nome: Ajouter les méthodes suivantes:
// - func CreateRoom(name, gameType string, hostID int) (*Room, error)
// - func GetRoomByCode(code string) (*Room, error)
// - func JoinRoom(roomID, userID int) error
// - func GetRoomParticipants(roomID int) ([]*RoomParticipant, error)
// - func GenerateRoomCode() string // Génère un code aléatoire de 6 caractères
