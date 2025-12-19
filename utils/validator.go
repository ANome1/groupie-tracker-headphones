package utils

import (
	"log"
	"regexp"
)

func ValidateEmail(email string) bool {
	re := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	log.Printf("Validating email: %s", email)
	return re.MatchString(email)
}

func ValidateUsername(username string) bool {
	re := regexp.MustCompile(`^[a-zA-Z0-9_]{3,20}$`)
	log.Printf("Validating username: %s", username)
	return re.MatchString(username)
}

// ValidatePassword valide le mot de passe selon les recommandations CNIL:
// - Au minimum 12 caractères OU
// - Au minimum 8 caractères avec complexité (majuscules + minuscules + chiffres + caractères spéciaux)
func ValidatePassword(password string) bool {
	log.Printf("Validating password: length=%d", len(password))

	// Règle 1: Au moins 12 caractères
	if len(password) >= 12 {
		return true
	}

	// Règle 2: Au minimum 8 caractères avec complexité
	if len(password) < 8 {
		return false
	}

	hasUppercase := regexp.MustCompile(`[A-Z]`).MatchString(password)
	hasLowercase := regexp.MustCompile(`[a-z]`).MatchString(password)
	hasDigit := regexp.MustCompile(`[0-9]`).MatchString(password)
	hasSpecial := regexp.MustCompile(`[!@#$%^&*()_+\-=\[\]{};:'"\\|,.<>\/?]`).MatchString(password)

	return hasUppercase && hasLowercase && hasDigit && hasSpecial
}
