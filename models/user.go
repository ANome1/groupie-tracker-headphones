package models

import "time"

// TODO @Nome: Structure User (ID, Username, Email, PasswordHash, CreatedAt)
// TODO @Nome: Fonctions CRUD
type User struct {
	ID           int
	Username     string
	Email        string
	PasswordHash string
	CreatedAt    time.Time
}
