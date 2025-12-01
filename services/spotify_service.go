package services

// RESPONSABLE: @Quoc Huy
// Service pour l'intégration avec l'API Spotify

// SpotifyService - Gestion de l'API Spotify
type SpotifyService struct {
	clientID     string
	clientSecret string
	accessToken  string
}

// NewSpotifyService - Crée une nouvelle instance du service
// TODO @Quoc Huy:
// - Charger les credentials depuis la config
func NewSpotifyService(clientID, clientSecret string) *SpotifyService {
	return &SpotifyService{
		clientID:     clientID,
		clientSecret: clientSecret,
	}
}

// Authenticate - Obtient un token d'accès Spotify
// TODO @Quoc Huy:
// - Utiliser le flow Client Credentials
// - Endpoint: https://accounts.spotify.com/api/token
// - Stocker l'access token
func (s *SpotifyService) Authenticate() error {
	// TODO: Implementation
	return nil
}

// GetRandomTracksFromPlaylist - Récupère des morceaux aléatoires
// TODO @Quoc Huy:
// - Paramètres: playlist ("Rock", "Rap", "Pop"), count (nombre de morceaux)
// - Utiliser l'API Spotify Search ou les playlists officielles:
//   - Rock: playlist_id ou search query "genre:rock"
//   - Rap: playlist_id ou search query "genre:rap"
//   - Pop: playlist_id ou search query "genre:pop"
//
// - Récupérer les preview_url (extraits de 30s)
// - Retourner une liste de Track (voir models/blindtest.go)
func (s *SpotifyService) GetRandomTracksFromPlaylist(playlist string, count int) ([]*Track, error) {
	// TODO: Implementation
	// Endpoint: https://api.spotify.com/v1/search?q=genre:rock&type=track&limit=50
	// Puis sélectionner 'count' morceaux aléatoirement
	return nil, nil
}

// SearchTrack - Recherche une musique spécifique
// TODO @Quoc Huy:
// - Utile pour valider les réponses des joueurs
// - Paramètres: query (titre + artiste)
func (s *SpotifyService) SearchTrack(query string) (*Track, error) {
	// TODO: Implementation
	return nil, nil
}

type Track struct {
	ID         string
	Title      string
	Artist     string
	PreviewURL string
	Album      string
}
