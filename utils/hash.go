package utils

import (
	"crypto/sha256"
	"encoding/hex"
	"log"
)

func HashPassword(password string) string {
	hash := sha256.Sum256([]byte(password))
	log.Printf("Hashed password: %x", hash)
	return hex.EncodeToString(hash[:])
}

func CheckPasswordHash(password, hash string) bool {
	return HashPassword(password) == hash
}
