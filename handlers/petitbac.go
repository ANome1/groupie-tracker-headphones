package handlers

// RESPONSABLE: @ilian
// Handlers pour le jeu Petit Bac

import "net/http"

// PetitBacGameHandler - Affiche l'interface du jeu Petit Bac
// TODO @ilian:
//   - Vérifier que le joueur est dans la salle
//   - Afficher l'interface de jeu (templates/games/petitbac.html):
//     → Grande lettre affichée au centre
//     → Grille de saisie avec les catégories
//     → Timer visuel pour la manche
//     → Scoreboard en temps réel (scoreboardActualPointInGame)
//     → Bouton pour gérer les catégories personnalisées
func PetitBacGameHandler(w http.ResponseWriter, r *http.Request) {
	// TODO: Implementation
}

// StartPetitBacHandler - Démarre une partie de Petit Bac
// TODO @ilian:
// - Vérifier que l'utilisateur est l'hôte de la salle
// - Récupérer les catégories sélectionnées
// - Initialiser le jeu avec NbrsManche = 9 manches
// - Changer le statut de la salle en 'in_progress'
// - Notifier tous les joueurs via WebSocket
// - Lancer le premier round (tirer une lettre)
func StartPetitBacHandler(w http.ResponseWriter, r *http.Request) {
	// TODO: Implementation
}

// SubmitPetitBacAnswersHandler - Soumet les réponses d'un joueur
// TODO @ilian:
// - Récupérer toutes les réponses du joueur (une par catégorie)
// - Vérifier que chaque réponse commence par la lettre du round
// - Stocker les réponses pour validation collective
// - Quand tous les joueurs ont soumis, lancer la validation
func SubmitPetitBacAnswersHandler(w http.ResponseWriter, r *http.Request) {
	// TODO: Implementation
}

// ValidateAnswersHandler - Gère la validation collective des réponses
// TODO @ilian:
//   - Afficher toutes les réponses de tous les joueurs
//   - Permettre à chaque joueur de voter pour chaque réponse
//   - SYSTÈME DE VALIDATION:
//     → Une réponse est validée si 2/3 des joueurs votent pour
//   - SYSTÈME DE POINTS:
//     → Réponse unique (personne d'autre): 2 points
//     → Réponse commune (au moins 2 joueurs): 1 point
//     → Réponse invalidée ou vide: 0 point
//   - Mettre à jour scoreboardActualPointInGame
//   - Notifier tous les joueurs des résultats via WebSocket
func ValidateAnswersHandler(w http.ResponseWriter, r *http.Request) {
	// TODO: Implementation
}

// NextPetitBacRoundHandler - Passe à la manche suivante
// TODO @ilian:
// - Vérifier que l'utilisateur est l'hôte
// - Tirer une nouvelle lettre aléatoire
// - Réinitialiser le timer
// - Notifier tous les joueurs via WebSocket
// - Si toutes les manches (NbrsManche) sont terminées, afficher le scoreboard final
func NextPetitBacRoundHandler(w http.ResponseWriter, r *http.Request) {
	// TODO: Implementation
}

// ============ GESTION DES CATÉGORIES PERSONNALISÉES ============

// GetCategoriesHandler - Liste toutes les catégories (défaut + personnalisées)
// TODO @ilian:
// - Récupérer toutes les catégories depuis la base de données
// - Retourner en JSON
func GetCategoriesHandler(w http.ResponseWriter, r *http.Request) {
	// TODO: Implementation
}

// CreateCategoryHandler - Crée une nouvelle catégorie personnalisée
// TODO @ilian:
// - Récupérer le nom de la catégorie
// - Vérifier qu'elle n'existe pas déjà
// - Créer en base de données
// - Retourner la catégorie créée
func CreateCategoryHandler(w http.ResponseWriter, r *http.Request) {
	// TODO: Implementation
}

// UpdateCategoryHandler - Modifie une catégorie existante
// TODO @ilian:
// - Vérifier que la catégorie n'est pas une catégorie par défaut
// - Mettre à jour le nom
func UpdateCategoryHandler(w http.ResponseWriter, r *http.Request) {
	// TODO: Implementation
}

// DeleteCategoryHandler - Supprime une catégorie personnalisée
// TODO @ilian:
// - Vérifier que la catégorie n'est pas une catégorie par défaut
// - Supprimer de la base de données
func DeleteCategoryHandler(w http.ResponseWriter, r *http.Request) {
	// TODO: Implementation
}
