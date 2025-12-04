package handlers

import (
	"groupie-tracker/utils"
	"html/template"
	"log"
	"net/http"
)

type RegisterData struct {
	Error string
}

// TODO @Nome: Handler POST pour inscription (récupérer form, valider, hasher SHA256, insérer en DB)
// TODO @Nome: Handler POST pour connexion (vérifier credentials, créer session)
// TODO @Nome: Handler pour déconnexion
func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		username := r.FormValue("username")
		email := r.FormValue("email")
		password := r.FormValue("password")

		if !utils.ValidateUsername(username) {
			tmpl, _ := template.ParseFiles("./templates/auth/register.html", "./templates/header.html", "./templates/footer.html")
			tmpl.Execute(w, RegisterData{Error: "Username invalide (3-20 caractères)"})
			return
		}

		if !utils.ValidateEmail(email) {
			tmpl, _ := template.ParseFiles("./templates/auth/register.html", "./templates/header.html", "./templates/footer.html")
			tmpl.Execute(w, RegisterData{Error: "Email invalide"})
			return
		}

		if !utils.ValidatePassword(password) {
			tmpl, _ := template.ParseFiles("./templates/auth/register.html", "./templates/header.html", "./templates/footer.html")
			tmpl.Execute(w, RegisterData{Error: "Mot de passe trop court (min 8 caractères)"})
			return
		}

		userID, err := AuthService.CreateUser(username, email, password)
		if err != nil {
			tmpl, _ := template.ParseFiles("./templates/auth/register.html", "./templates/header.html", "./templates/footer.html")
			tmpl.Execute(w, RegisterData{Error: "Erreur : " + err.Error()})
			return
		}

		log.Printf("Utilisateur créé : %d (%s)", userID, username)
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	tmpl, err := template.ParseFiles("./templates/auth/register.html", "./templates/header.html", "./templates/footer.html")
	if err != nil {
		log.Printf("Erreur: %v", err)
		http.Error(w, "Erreur serveur", http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, RegisterData{Error: ""})
}
