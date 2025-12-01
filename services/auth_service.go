package services

// RESPONSABLE: @Nome
// Service pour la logique métier de l'authentification

// AuthService - Gestion de l'authentification
type AuthService struct {
	// TODO @Nome: Ajouter la connexion à la base de données
}

// NewAuthService - Crée une nouvelle instance du service
// TODO @Nome:
func NewAuthService() *AuthService {
	return &AuthService{}
}

// Register - Inscrit un nouvel utilisateur
// TODO @Nome:
// - Valider les données (utils/validator.go)
// - Vérifier que l'username et l'email sont uniques
// - Hasher le mot de passe (utils/hash.go)
// - Insérer en base de données
// - Retourner l'ID de l'utilisateur créé
func (s *AuthService) Register(username, email, password string) (int, error) {
	// TODO: Implementation
	return 0, nil
}

// Login - Authentifie un utilisateur
// TODO @Nome:
// - Récupérer l'utilisateur par username
// - Vérifier le mot de passe (utils/hash.go)
// - Retourner l'ID de l'utilisateur
func (s *AuthService) Login(username, password string) (int, error) {
	// TODO: Implementation
	return 0, nil
}

// GetUserByID - Récupère un utilisateur par son ID
// TODO @Nome:
func (s *AuthService) GetUserByID(userID int) (*User, error) {
	// TODO: Implementation
	return nil, nil
}

type User struct {
	ID       int
	Username string
	Email    string
}
