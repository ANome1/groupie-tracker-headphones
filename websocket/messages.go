package websocket

// RESPONSABLE: @Quoc Huy & @ilian
// Définition des types de messages WebSocket

// MessageType - Type de message
type MessageType string

const (
	// Messages généraux
	MessageTypeJoinRoom  MessageType = "join_room"
	MessageTypeLeaveRoom MessageType = "leave_room"
	MessageTypeRoomState MessageType = "room_state"

	// Messages de jeu
	MessageTypeGameStart MessageType = "game_start"
	MessageTypeGameEnd   MessageType = "game_end"
	MessageTypeNextRound MessageType = "next_round"

	// TODO @Quoc Huy: Messages Blind Test
	MessageTypeBlindTestTrack  MessageType = "blindtest_track"  // Nouvelle musique
	MessageTypeBlindTestAnswer MessageType = "blindtest_answer" // Réponse d'un joueur

	// TODO @ilian: Messages Petit Bac
	MessageTypePetitBacLetter   MessageType = "petitbac_letter"   // Nouvelle lettre
	MessageTypePetitBacAnswers  MessageType = "petitbac_answers"  // Réponses d'un joueur
	MessageTypePetitBacValidate MessageType = "petitbac_validate" // Validation collective

	// Messages de score
	MessageTypeScoreUpdate MessageType = "score_update"

	// Erreurs
	MessageTypeError MessageType = "error"
)

// Message - Structure générique d'un message WebSocket
type Message struct {
	Type    MessageType            `json:"type"`
	RoomID  int                    `json:"room_id,omitempty"`
	UserID  int                    `json:"user_id,omitempty"`
	Payload map[string]interface{} `json:"payload,omitempty"`
}

// TODO @Quoc Huy: Messages spécifiques au Blind Test
// BlindTestTrackMessage - Envoie une nouvelle musique
type BlindTestTrackMessage struct {
	Type       MessageType `json:"type"`
	RoomID     int         `json:"room_id"`
	TrackID    string      `json:"track_id"`
	PreviewURL string      `json:"preview_url"`
	Timer      int         `json:"timer"` // 37 secondes
	Round      int         `json:"round"`
}

// BlindTestAnswerMessage - Réponse d'un joueur
type BlindTestAnswerMessage struct {
	Type        MessageType `json:"type"`
	RoomID      int         `json:"room_id"`
	UserID      int         `json:"user_id"`
	Answer      string      `json:"answer"`
	TimeElapsed int         `json:"time_elapsed"` // Temps depuis le début de la musique
}

// TODO @ilian: Messages spécifiques au Petit Bac
// PetitBacLetterMessage - Envoie une nouvelle lettre
type PetitBacLetterMessage struct {
	Type   MessageType `json:"type"`
	RoomID int         `json:"room_id"`
	Letter string      `json:"letter"`
	Timer  int         `json:"timer"`
	Round  int         `json:"round"`
}

// PetitBacAnswersMessage - Réponses d'un joueur
type PetitBacAnswersMessage struct {
	Type    MessageType       `json:"type"`
	RoomID  int               `json:"room_id"`
	UserID  int               `json:"user_id"`
	Answers map[string]string `json:"answers"` // category -> answer
}

// PetitBacValidateMessage - Vote de validation
type PetitBacValidateMessage struct {
	Type   MessageType             `json:"type"`
	RoomID int                     `json:"room_id"`
	UserID int                     `json:"user_id"`
	Votes  map[int]map[string]bool `json:"votes"` // targetUserID -> category -> valid
}

// ScoreUpdateMessage - Mise à jour des scores
type ScoreUpdateMessage struct {
	Type   MessageType `json:"type"`
	RoomID int         `json:"room_id"`
	Scores map[int]int `json:"scores"` // userID -> score
}

// ErrorMessage - Message d'erreur
type ErrorMessage struct {
	Type    MessageType `json:"type"`
	Message string      `json:"message"`
}
