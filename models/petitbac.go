package models

// RESPONSABLE: @ilian
// Modèle pour le jeu Petit Bac

// NOMBRE DE MANCHES PAR DÉFAUT
// Cette constante définit le nombre total de manches dans une partie de Petit Bac
// Chaque manche correspond à une nouvelle lettre tirée aléatoirement
const NbrsManche = 9

type PetitBacGame struct {
	CurrentRound  int              `json:"current_round"`
	TotalRounds   int              `json:"total_rounds"`   // Utiliser NbrsManche
	CurrentLetter string           `json:"current_letter"` // Lettre tirée (A-Z)
	Categories    []string         `json:"categories"`     // Liste des catégories
	PlayerAnswers map[int][]Answer `json:"player_answers"` // userID -> réponses
	PlayerScores  map[int]int      `json:"player_scores"`  // userID -> score total
	RoundStarted  bool             `json:"round_started"`
	TimeRemaining int              `json:"time_remaining"` // En secondes

	// Variable pour suivre les points actuels dans la partie
	// Utilisé pour le scoreboard en temps réel pendant le jeu
	ScoreboardActualPointInGame map[int]int `json:"scoreboard_actual_point_in_game"`
}

type Answer struct {
	Category  string `json:"category"`
	Value     string `json:"value"`
	Points    int    `json:"points"`    // 0, 1, ou 2 points
	Validated bool   `json:"validated"` // Validé par les autres joueurs
}

type Category struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	CreatedBy int    `json:"created_by,omitempty"`
	IsDefault bool   `json:"is_default"`
}

func (g *PetitBacGame) GetType() string {
	return "petitbac"
}

func (g *PetitBacGame) IsFinished() bool {
	return g.CurrentRound >= g.TotalRounds
}

func (g *PetitBacGame) GetScores() map[int]int {
	return g.PlayerScores
}

// TODO @ilian: Implémenter les fonctions suivantes:
// - func NewPetitBacGame(categories []string) *PetitBacGame
//   → Initialiser avec NbrsManche = 9 manches
//   → Initialiser ScoreboardActualPointInGame
//
// - func (g *PetitBacGame) StartRound()
//   → Tirer une lettre aléatoire (A-Z, éviter K, W, X, Y, Z si trop difficile)
//
// - func (g *PetitBacGame) SubmitAnswers(userID int, answers []Answer)
//   → Stocker les réponses du joueur pour validation
//
// - func (g *PetitBacGame) ValidateAnswers() map[int]map[string]int
//   → SYSTÈME DE VALIDATION COLLECTIVE:
//     * Chaque réponse est votée par les autres joueurs
//     * Si 2/3 des joueurs valident → réponse acceptée
//   → SYSTÈME DE POINTS:
//     * Réponse unique (personne d'autre n'a la même): 2 points
//     * Réponse commune (d'autres joueurs ont la même): 1 point
//     * Réponse invalidée ou vide: 0 point
//   → Mettre à jour ScoreboardActualPointInGame
//
// - func (g *PetitBacGame) NextRound()
//   → Passer à la manche suivante
//
// - func (g *PetitBacGame) GetLeaderboard() []struct{UserID int; Score int}

// TODO @ilian: Gestion des catégories personnalisées (CRUD)
// - func GetAllCategories() ([]*Category, error)
// - func CreateCategory(name string, userID int) (*Category, error)
// - func UpdateCategory(id int, name string) error
// - func DeleteCategory(id int) error
// - func GetDefaultCategories() ([]*Category, error)
//   → Catégories par défaut: Prénom, Ville, Pays, Animal, Métier, etc.
