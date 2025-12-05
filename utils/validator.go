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

func ValidatePassword(password string) bool {
	log.Printf("Validating password length: %d", len(password))
	return len(password) >= 8
}
