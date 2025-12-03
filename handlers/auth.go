package handlers

// TODO @Nome: Handler POST pour inscription (récupérer form, valider, hasher SHA256, insérer en DB)
// TODO @Nome: Handler POST pour connexion (vérifier credentials, créer session)
// TODO @Nome: Handler pour déconnexion
import (
	"crypto/sha256"
	"encoding/hex"
	"log"
	"regexp"
)

func HashPassword(password string) string {
	hash := sha256.Sum256([]byte(password))
	log.Printf("Hashed password: %x", hash)
	return hex.EncodeToString(hash[:])
}
func CheckPasswordHash(password, hash string) bool {
	return HashPassword(password) == hash
}

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
