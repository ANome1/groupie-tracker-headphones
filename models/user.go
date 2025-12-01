package models

// RESPONSABLE: @Nome
// Modèle utilisateur pour l'authentification

import "time"

type User struct {
	ID           int       `json:"id"`
	Username     string    `json:"username"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"-"` // Ne pas exposer le hash dans le JSON
	CreatedAt    time.Time `json:"created_at"`
}

// TODO @Nome: Ajouter les méthodes suivantes:
// - func CreateUser(username, email, passwordHash string) (*User, error)
// - func GetUserByUsername(username string) (*User, error)
// - func GetUserByEmail(email string) (*User, error)
// - func GetUserByID(id int) (*User, error)
