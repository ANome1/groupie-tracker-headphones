package models

// RESPONSABLE: @Quoc Huy & @ilian
// Interface commune pour les jeux

// TODO @Quoc Huy & @ilian: Créer la structure GameSession avec:
// - ID int
// - RoomID int
// - GameData string (JSON de l'état du jeu)
// - StartedAt time.Time
// - EndedAt *time.Time

// TODO @Quoc Huy & @ilian: Créer l'interface Game avec:
// - GetType() string (retourne "blindtest" ou "petitbac")
// - IsFinished() bool (vérifie si le jeu est terminé)
// - GetScores() map[int]int (retourne userID -> score)
