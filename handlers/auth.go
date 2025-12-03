package handlers

// TODO @Nome: Handler POST pour inscription (récupérer form, valider, hasher SHA256, insérer en DB)
// TODO @Nome: Handler POST pour connexion (vérifier credentials, créer session)
// TODO @Nome: Handler pour déconnexion
import (
	"crypto/sha256"
	"encoding/hex"
	"regexp"
)

func HashPassword(password string) string {
	hash := sha256.Sum256([]byte(password))
	return hex.EncodeToString(hash[:])
}
func CheckPasswordHash(password, hash string) bool {
	return HashPassword(password) == hash
}

func ValidateEmail(email string) bool {
	re := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return re.MatchString(email)
}

func ValidateUsername(username string) bool {
	re := regexp.MustCompile(`^[a-zA-Z0-9_]{3,20}$`)
	return re.MatchString(username)
}

func ValidatePassword(password string) bool {
	return len(password) >= 8
}
