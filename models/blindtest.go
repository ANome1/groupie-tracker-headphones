package models

// RESPONSABLE: @Quoc Huy
// Modèle pour le jeu Blind Test

// TODO @Quoc Huy: Définir la constante DefaultBlindTestTimer = 37 secondes

// TODO @Quoc Huy: Créer la structure BlindTestGame avec:
// - RoomID int
// - CurrentTrack *Track
// - Tracks []*Track
// - CurrentIndex int
// - StartTime time.Time
// - PlayerAnswers map[int]*Answer (userID -> Answer)
// - Scores map[int]int (userID -> score)

// TODO @Quoc Huy: Créer la structure Track avec:
// - SpotifyID string
// - Title string
// - Artist string
// - PreviewURL string (URL audio 30s)

// TODO @Quoc Huy: Créer la structure Answer avec:
// - UserID int
// - Title string
// - Artist string
// - Timestamp time.Time (moment de la réponse)
// - IsCorrect bool
// - Points int (3, 2, 1 ou 0 points)

// TODO @Quoc Huy: Implémenter les méthodes de l'interface Game:
// - func (g *BlindTestGame) GetType() string { return "blindtest" }
// - func (g *BlindTestGame) IsFinished() bool
// - func (g *BlindTestGame) GetScores() map[int]int

// TODO @Quoc Huy: Système de points (selon le temps de réponse):
// - Réponse correcte en moins de 10s: 3 points
// - Réponse correcte en moins de 20s: 2 points
// - Réponse correcte en moins de 37s: 1 point
// - Pas de réponse ou fausse réponse: 0 point

// TODO @Quoc Huy: Ajouter une fonction pour valider les réponses:
// - func ValidateAnswer(userAnswer, correctAnswer string) bool
// - Ignorer la casse et les espaces supplémentaires
// - Accepter les réponses partielles (ex: "Bohemian" pour "Bohemian Rhapsody")

// TODO @Quoc Huy: Intégration Spotify API
// - func GetRandomTracksFromPlaylist(playlist string, count int) ([]*Track, error)
//   → Utiliser l'API Spotify pour récupérer des morceaux aléatoires
//   → Playlists disponibles: "Rock", "Rap", "Pop"
