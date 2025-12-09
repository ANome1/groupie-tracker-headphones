package handlers

import (
	"groupie-tracker/utils"
	"html/template"
	"log"
	"net/http"
	// "strconv" // TODO: Décommentez quand les cookies seront réactivés
)

type RegisterData struct {
	Error string
}

type LoginData struct {
	Error string
}

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

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		username := r.FormValue("username")
		password := r.FormValue("password")

		if username == "" || password == "" {
			tmpl, _ := template.ParseFiles("./templates/auth/login.html", "./templates/header.html", "./templates/footer.html")
			tmpl.Execute(w, LoginData{Error: "Nom d'utilisateur et mot de passe requis"})
			return
		}

		user, err := AuthService.ValidateUserCredentials(username, password)
		if err != nil {
			tmpl, _ := template.ParseFiles("./templates/auth/login.html", "./templates/header.html", "./templates/footer.html")
			tmpl.Execute(w, LoginData{Error: "Identifiants invalides"})
			return
		}

		// TODO: Cookies désactivés pour développement du front
		/*
			http.SetCookie(w, &http.Cookie{
				Name:   "user_id",
				Value:  strconv.Itoa(user.ID),
				MaxAge: 3600 * 24 * 7,
				Path:   "/",
			})
		*/

		log.Printf("Utilisateur connecté : %s (ID: %d)", user.Username, user.ID)
		http.Redirect(w, r, "/home", http.StatusSeeOther)
		return
	}

	tmpl, err := template.ParseFiles("./templates/auth/login.html", "./templates/header.html", "./templates/footer.html")
	if err != nil {
		log.Printf("Erreur: %v", err)
		http.Error(w, "Erreur serveur", http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, LoginData{Error: ""})
}

func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	// TODO: Cookies désactivés pour développement du front
	/*
		http.SetCookie(w, &http.Cookie{
			Name:   "user_id",
			MaxAge: -1,
			Path:   "/",
		})
	*/

	log.Printf("Utilisateur déconnecté")
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
