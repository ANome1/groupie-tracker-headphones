package config

import (
	// "database/sql"
	_ "github.com/mattn/go-sqlite3"
)

type Config struct {
	DatabasePath string
	ServerPort   string
}

func Load() *Config {
	return &Config{
		DatabasePath: "",
		ServerPort:   "8080",
	}
}
