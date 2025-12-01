package handlers

// RESPONSABLE: @ilian
// Handlers pour le jeu Petit Bac

// TODO @ilian: PetitBacHandler - Page principale du Petit Bac
// - Afficher l'interface de jeu (templates/games/petitbac.html)
// - Charger les informations de la salle et les catégories
// - Vérifier que le joueur appartient à cette salle

// TODO @ilian: SubmitAnswersHandler - Soumettre les réponses
// - Méthode POST
// - Récupérer les réponses (une par catégorie)
// - Valider que toutes commencent par la bonne lettre
// - Stocker en attente de validation collective
// - Déclencher la phase de validation via WebSocket

// TODO @ilian: ValidateAnswersHandler - Voter pour valider les réponses
// - Méthode POST
// - Récupérer les votes (valid/invalid pour chaque réponse)
// - SYSTÈME DE VALIDATION COLLECTIVE:
//   * Une réponse est acceptée si 2/3 des joueurs la valident
//   * Ex: 3 joueurs = 2 votes, 4 joueurs = 3 votes
// - Calculer les points:
//   * Réponse unique: 2 points
//   * Réponse commune: 1 point
//   * Réponse invalidée: 0 point
// - Mettre à jour scoreboardActualPointInGame
// - Envoyer les résultats via WebSocket

// TODO @ilian: CRUD Catégories personnalisées
// - CreateCategoryHandler: Créer une nouvelle catégorie
// - ListCategoriesHandler: Lister toutes les catégories
// - UpdateCategoryHandler: Modifier une catégorie
// - DeleteCategoryHandler: Supprimer une catégorie
