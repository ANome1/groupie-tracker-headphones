package services

// RESPONSABLE: @Quoc Huy & @ilian
// Service pour la logique commune des jeux

// GameService - Gestion des sessions de jeu
type GameService struct {
	// TODO: Ajouter la connexion à la base de données
}

// NewGameService - Crée une nouvelle instance du service
func NewGameService() *GameService {
	return &GameService{}
}

// CreateGameSession - Crée une nouvelle session de jeu
// TODO @Quoc Huy & @ilian:
// - Insérer une entrée dans la table game_sessions
// - Retourner l'ID de la session
func (s *GameService) CreateGameSession(roomID int, gameType string) (int, error) {
	// TODO: Implementation
	return 0, nil
}

// UpdateGameState - Sauvegarde l'état actuel du jeu
// TODO @Quoc Huy & @ilian:
// - Sérialiser l'état du jeu en JSON
// - Mettre à jour la colonne game_data
func (s *GameService) UpdateGameState(sessionID int, gameData interface{}) error {
	// TODO: Implementation
	return nil
}

// EndGameSession - Termine une session de jeu
// TODO @Quoc Huy & @ilian:
// - Mettre à jour ended_at
// - Sauvegarder les scores finaux
func (s *GameService) EndGameSession(sessionID int) error {
	// TODO: Implementation
	return nil
}
