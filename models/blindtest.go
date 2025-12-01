package models

// RESPONSABLE: @Quoc Huy
// Modèle pour le jeu Blind Test

type BlindTestGame struct {
	CurrentRound  int         `json:"current_round"`
	TotalRounds   int         `json:"total_rounds"`
	CurrentTrack  *Track      `json:"current_track,omitempty"`
	PlayerScores  map[int]int `json:"player_scores"` // userID -> score
	RoundStarted  bool        `json:"round_started"`
	TimeRemaining int         `json:"time_remaining"` // En secondes (défaut: 37)
	Playlist      string      `json:"playlist"`       // "Rock", "Rap", ou "Pop"
}

// TIMER PAR DÉFAUT: 37 secondes par musique
const DefaultBlindTestTimer = 37

type Track struct {
	ID         string `json:"id"`          // ID Spotify
	Title      string `json:"title"`       // Titre de la chanson
	Artist     string `json:"artist"`      // Artiste
	PreviewURL string `json:"preview_url"` // URL de l'extrait (30s)
	Album      string `json:"album,omitempty"`
	Year       int    `json:"year,omitempty"`
}

func (g *BlindTestGame) GetType() string {
	return "blindtest"
}

func (g *BlindTestGame) IsFinished() bool {
	return g.CurrentRound >= g.TotalRounds
}

func (g *BlindTestGame) GetScores() map[int]int {
	return g.PlayerScores
}

// TODO @Quoc Huy: Implémenter les fonctions suivantes:
// - func NewBlindTestGame(playlist string, totalRounds int) *BlindTestGame
// - func (g *BlindTestGame) StartRound(track *Track)
// - func (g *BlindTestGame) SubmitAnswer(userID int, answer string, timeElapsed int) (correct bool, points int)
//   → Système de points basé sur la rapidité:
//     * Réponse correcte en < 10s: 3 points
//     * Réponse correcte en < 20s: 2 points
//     * Réponse correcte en < 37s: 1 point
// - func (g *BlindTestGame) NextRound()
// - func (g *BlindTestGame) GetLeaderboard() []struct{UserID int; Score int}

// TODO @Quoc Huy: Intégration Spotify API
// - func GetRandomTracksFromPlaylist(playlist string, count int) ([]*Track, error)
//   → Utiliser l'API Spotify pour récupérer des morceaux aléatoires
//   → Playlists disponibles: "Rock", "Rap", "Pop"
