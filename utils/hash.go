package utils

// RESPONSABLE: @Nome
// Fonctions de hashing des mots de passe
// IMPORTANT: Ne pas utiliser golang.org/x/crypto/bcrypt
// Utiliser crypto/sha256 de la bibliothèque standard

import (
	"crypto/sha256"
	"encoding/hex"
)

// HashPassword - Hash un mot de passe
// TODO @Nome:
// - Utiliser SHA256
// - Ajouter un salt (sel) pour plus de sécurité
// - Retourner le hash en hexadécimal
func HashPassword(password string) string {
	// Exemple simple (AMÉLIORER avec un salt unique par utilisateur)
	hash := sha256.Sum256([]byte(password))
	return hex.EncodeToString(hash[:])
}

// CheckPasswordHash - Vérifie qu'un mot de passe correspond au hash
// TODO @Nome:
// - Hasher le mot de passe fourni
// - Comparer avec le hash stocké
// - Retourner true si identique, false sinon
func CheckPasswordHash(password, hash string) bool {
	passwordHash := HashPassword(password)
	return passwordHash == hash
}

// TODO @Nome: AMÉLIORATION - Implémenter un salt unique
// - Stocker le salt avec le hash dans la base de données
// - Format: "salt$hash"
// func HashPasswordWithSalt(password string) string
// func CheckPasswordHashWithSalt(password, storedHash string) bool
