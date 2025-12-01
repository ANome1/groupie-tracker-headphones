package utils

// RESPONSABLE: @Nome, @Quoc Huy, @ilian
// Fonctions utilitaires partagées

import (
	"encoding/json"
	"net/http"
)

// RespondJSON - Envoie une réponse JSON
// TODO @Nome:
// - Définir le Content-Type
// - Encoder les données en JSON
// - Gérer les erreurs d'encodage
func RespondJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

// RespondError - Envoie une erreur en JSON
// TODO @Nome:
func RespondError(w http.ResponseWriter, status int, message string) {
	RespondJSON(w, status, map[string]string{"error": message})
}

// ParseJSON - Parse le corps d'une requête JSON
// TODO @Nome:
func ParseJSON(r *http.Request, v interface{}) error {
	return json.NewDecoder(r.Body).Decode(v)
}

// TODO @Quoc Huy: Fonction pour comparer deux strings (titre/artiste)
// - Ignorer la casse
// - Ignorer les accents
// - Ignorer la ponctuation
// func NormalizeString(s string) string
// func CompareStrings(s1, s2 string) bool

// TODO @ilian: Fonction pour tirer une lettre aléatoire
// - Exclure les lettres difficiles (K, W, X, Y, Z) optionnel
// func RandomLetter() string
