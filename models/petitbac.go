package models

// RESPONSABLE: @ilian
// Modèle pour le jeu Petit Bac

// TODO @ilian: Définir la constante NbrsManche = 9 (nombre de manches par partie)

// TODO @ilian: Créer la structure PetitBacGame avec:
// - RoomID int
// - Categories []string (ex: ["Pays", "Prénom", "Animal", ...])
// - CurrentLetter string
// - CurrentRound int
// - PlayerAnswers map[int]map[string]string (userID -> catégorie -> réponse)
// - ValidationVotes map[int]map[string][]int (userID -> catégorie -> liste voters)
// - Scores map[int]int (userID -> score total)
// - scoreboardActualPointInGame map[int]int (userID -> points de cette manche)
// - StartTime time.Time

// TODO @ilian: Créer la structure PetitBacAnswer avec:
// - UserID int
// - Category string
// - Answer string
// - IsValid bool (validé par le vote collectif)
// - IsUnique bool (unique = 2 points, commun = 1 point)
// - Points int

// TODO @ilian: Implémenter les méthodes de l'interface Game:
// - func (g *PetitBacGame) GetType() string { return "petitbac" }
// - func (g *PetitBacGame) IsFinished() bool { return g.CurrentRound >= NbrsManche }
// - func (g *PetitBacGame) GetScores() map[int]int

// TODO @ilian: Système de validation collective:
// - Chaque réponse doit être validée par au moins 2/3 des joueurs
// - Exemple: 3 joueurs = 2 votes requis, 4 joueurs = 3 votes requis
// - Afficher les réponses de tous les joueurs pour la phase de vote

// TODO @ilian: Système de points:
// - Réponse unique (personne d'autre n'a donné la même): 2 points
// - Réponse commune (au moins 2 joueurs ont la même): 1 point
// - Réponse invalide ou vide: 0 point

// TODO @ilian: Gestion des catégories personnalisées (CRUD):
// - Les joueurs peuvent créer leurs propres catégories
// - Stocker les catégories custom en base de données
// - L'hôte de la salle choisit les catégories avant de lancer la partie
