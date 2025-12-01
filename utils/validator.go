package utils

// RESPONSABLE: @Nome
// Fonctions de validation des données utilisateur

import "regexp"

// ValidateEmail - Vérifie qu'un email est valide
// TODO @Nome:
// - Vérifier le format avec une regex
// - Retourner true si valide, false sinon
func ValidateEmail(email string) bool {
	// Regex simple pour email
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)
	return emailRegex.MatchString(email)
}

// ValidateUsername - Vérifie qu'un username est valide
// TODO @Nome:
// - 3 à 20 caractères
// - Alphanumérique + underscores uniquement
// - Retourner true si valide, false sinon
func ValidateUsername(username string) bool {
	if len(username) < 3 || len(username) > 20 {
		return false
	}
	usernameRegex := regexp.MustCompile(`^[a-zA-Z0-9_]+$`)
	return usernameRegex.MatchString(username)
}

// ValidatePassword - Vérifie qu'un mot de passe est valide
// TODO @Nome:
// - Minimum 8 caractères
// - Retourner true si valide, false sinon
func ValidatePassword(password string) bool {
	return len(password) >= 8
}

// ValidateRoomName - Vérifie qu'un nom de salle est valide
// TODO @Nome:
// - 3 à 50 caractères
func ValidateRoomName(name string) bool {
	return len(name) >= 3 && len(name) <= 50
}

// ValidateRoomCode - Vérifie qu'un code de salle est valide
// TODO @Nome:
// - 6 caractères exactement
// - Lettres majuscules et chiffres uniquement
func ValidateRoomCode(code string) bool {
	if len(code) != 6 {
		return false
	}
	codeRegex := regexp.MustCompile(`^[A-Z0-9]{6}$`)
	return codeRegex.MatchString(code)
}
