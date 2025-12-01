package config

// RESPONSABLE: @Nome
// Configuration de l'application

import "os"

// Config - Structure de configuration
type Config struct {
	DatabasePath  string // Chemin vers la base SQLite
	ServerPort    string // Port du serveur (défaut: 8080)
	SpotifyAPIKey string // Clé API Spotify (@Quoc Huy)
	SpotifySecret string // Secret Spotify (@Quoc Huy)
}

// Load - Charge la configuration
// TODO @Nome:
// - Lire les variables d'environnement
// - Définir des valeurs par défaut
// - Retourner la config
func Load() *Config {
	return &Config{
		DatabasePath:  getEnv("DATABASE_PATH", "./database/groupie-tracker.db"),
		ServerPort:    getEnv("SERVER_PORT", "8080"),
		SpotifyAPIKey: getEnv("SPOTIFY_API_KEY", ""),
		SpotifySecret: getEnv("SPOTIFY_SECRET", ""),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// TODO @Quoc Huy: Ajouter les credentials Spotify
// - Créer un compte développeur Spotify
// - Obtenir Client ID et Client Secret
// - Les stocker dans un fichier .env (ne pas commit)
