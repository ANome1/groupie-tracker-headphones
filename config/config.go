package config

import (
	"os"
)

// RESPONSABLE: @Nome
// Configuration de l'application

type Config struct {
	DatabasePath string
	ServerPort   string
}

func Load() *Config {
	dbPath := os.Getenv("DB_PATH")
	if dbPath == "" {
		dbPath = "./database/groupie-tracker.db"
	}

	serverPort := os.Getenv("SERVER_PORT")
	if serverPort == "" {
		serverPort = "8080"
	}

	return &Config{
		DatabasePath: dbPath,
		ServerPort:   serverPort,
	}
}
