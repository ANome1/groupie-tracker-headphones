package utils

import (
	"crypto/rand"
	"encoding/base64"
	"log"
	"net/http"
)

const CSRFTokenLength = 32
const CSRFCookieName = "csrf_token"

// GenerateCSRFToken génère un token CSRF aléatoire
func GenerateCSRFToken() string {
	b := make([]byte, CSRFTokenLength)
	_, err := rand.Read(b)
	if err != nil {
		log.Printf("Error generating CSRF token: %v", err)
		return ""
	}
	return base64.StdEncoding.EncodeToString(b)
}

// SetCSRFToken définit le token CSRF dans un cookie
func SetCSRFToken(w http.ResponseWriter, token string, isProduction bool) {
	http.SetCookie(w, &http.Cookie{
		Name:     CSRFCookieName,
		Value:    token,
		MaxAge:   3600, // 1 heure
		Path:     "/",
		HttpOnly: false, // Doit être accessible au JS pour les forms
		Secure:   isProduction,
		SameSite: http.SameSiteLaxMode,
	})
}

// GetCSRFToken récupère le token CSRF du cookie
func GetCSRFToken(r *http.Request) string {
	cookie, err := r.Cookie(CSRFCookieName)
	if err != nil {
		return ""
	}
	return cookie.Value
}

// VerifyCSRFToken vérifie que le token POST correspond au token du cookie
func VerifyCSRFToken(r *http.Request) bool {
	cookieToken := GetCSRFToken(r)
	if cookieToken == "" {
		log.Printf("CSRF: No token in cookie")
		return false
	}

	// Récupère le token du formulaire ou de l'en-tête
	formToken := r.FormValue("csrf_token")
	if formToken == "" {
		formToken = r.Header.Get("X-CSRF-Token")
	}

	if formToken == "" {
		log.Printf("CSRF: No token in request")
		return false
	}

	isValid := cookieToken == formToken
	if !isValid {
		log.Printf("CSRF: Token mismatch (cookie: %s, form: %s)", cookieToken[:10], formToken[:10])
	}
	return isValid
}
