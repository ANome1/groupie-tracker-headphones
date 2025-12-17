package services

import (
	"errors"
	"log"

	"groupie-tracker/database"
	"groupie-tracker/models"
	"groupie-tracker/utils"
)

type AuthService struct {
	DB *database.Database
}

func (as *AuthService) CreateUser(username, email, password string) (int, error) {

	if !utils.ValidateUsername(username) {
		return 0, errors.New("invalid username")
	}
	if !utils.ValidateEmail(email) {
		return 0, errors.New("invalid email")
	}
	if !utils.ValidatePassword(password) {
		return 0, errors.New("invalid password")
	}

	hashedPassword := utils.HashPassword(password)

	result, err := as.DB.DB.Exec("INSERT INTO users (username, email, password_hash) VALUES (?, ?, ?)", username, email, hashedPassword)
	if err != nil {
		log.Printf("Error creating user: %v", err)
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), nil
}

func (as *AuthService) GetUserByUsername(username string) (*models.User, error) {
	var user models.User
	err := as.DB.DB.QueryRow("SELECT id, username, email, password_hash, created_at FROM users WHERE username = ?", username).
		Scan(&user.ID, &user.Username, &user.Email, &user.PasswordHash, &user.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (as *AuthService) GetUserByID(id int) (*models.User, error) {
	var user models.User
	err := as.DB.DB.QueryRow("SELECT id, username, email, password_hash, created_at FROM users WHERE id = ?", id).
		Scan(&user.ID, &user.Username, &user.Email, &user.PasswordHash, &user.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (as *AuthService) ValidateUserCredentials(username, password string) (*models.User, error) {
	user, err := as.GetUserByUsername(username)
	if err != nil {
		return nil, err
	}

	if !utils.CheckPasswordHash(password, user.PasswordHash) {
		return nil, errors.New("invalid credentials")
	}

	return user, nil
}
