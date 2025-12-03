package database

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

type Database struct {
	DB *sql.DB
}

func InitDB(dbPath string) *Database {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		log.Fatalf("Erreur lors de l'ouverture de la base de données: %v", err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatalf("Erreur lors de la connexion à la base de données: %v", err)
	}

	schemaFile, err := os.ReadFile("./database/schema.sql")
	if err != nil {
		log.Fatalf("Erreur lors de la lecture du fichier schema.sql: %v", err)
	}

	schema := string(schemaFile)
	_, err = db.Exec(schema)
	if err != nil {
		log.Fatalf("Erreur lors de l'exécution du schéma de la base de données: %v", err)
	}
	return &Database{DB: db}
}

func (database *Database) Close() {
	err := database.DB.Close()
	if err != nil {
		log.Printf("Erreur lors de la fermeture de la base de données: %v", err)
	}
}
