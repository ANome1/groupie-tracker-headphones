package main

import (
	"groupie-tracker/config"
	"groupie-tracker/database"
	"groupie-tracker/handlers"
	"groupie-tracker/services"
	"log"
	"net/http"
	"text/template"
)

// RESPONSABLE: @Nome (infrastructure), @Quoc Huy (WebSocket Blind Test), @ilian (WebSocket Petit Bac)

// Landing - Page d'accueil
func Landing(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	tmpl, err := template.ParseFiles("./templates/landing.html", "./templates/header.html", "./templates/footer.html")
	if err != nil {
		log.Printf("Erreur: %v", err)
		http.Error(w, "Erreur serveur", http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, nil)
}

// Home - Page de sélection de jeu
func Home(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("./templates/home.html", "./templates/header.html", "./templates/footer.html")
	if err != nil {
		log.Printf("Erreur: %v", err)
		http.Error(w, "Erreur serveur", http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, nil)
}

// RegisterHandler - Inscription
// TODO @Nome: Récupérer username, email, password du formulaire
// TODO @Nome: Valider les données, hasher le mot de passe (SHA256)
// TODO @Nome: Insérer en base de données, créer une session
func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {

		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	tmpl, err := template.ParseFiles("./templates/auth/register.html", "./templates/header.html", "./templates/footer.html")
	if err != nil {
		log.Printf("Erreur: %v", err)
		http.Error(w, "Erreur serveur", http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, nil)
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		// TODO: Logique de connexion
		http.Redirect(w, r, "/home", http.StatusSeeOther)
		return
	}

	tmpl, err := template.ParseFiles("./templates/auth/login.html", "./templates/header.html", "./templates/footer.html")
	if err != nil {
		log.Printf("Erreur: %v", err)
		http.Error(w, "Erreur serveur", http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, nil)
}

func CreateRoomHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		http.Redirect(w, r, "/room/lobby", http.StatusSeeOther)
		return
	}

	tmpl, err := template.ParseFiles("./templates/room/create.html", "./templates/header.html", "./templates/footer.html")
	if err != nil {
		log.Printf("Erreur: %v", err)
		http.Error(w, "Erreur serveur", http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, nil)
}

func JoinRoomHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		http.Redirect(w, r, "/room/lobby", http.StatusSeeOther)
		return
	}

	tmpl, err := template.ParseFiles("./templates/room/join.html", "./templates/header.html", "./templates/footer.html")
	if err != nil {
		log.Printf("Erreur: %v", err)
		http.Error(w, "Erreur serveur", http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, nil)
}

func LobbyHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("./templates/room/lobby.html", "./templates/header.html", "./templates/footer.html")
	if err != nil {
		log.Printf("Erreur: %v", err)
		http.Error(w, "Erreur serveur", http.StatusInternalServerError)
		return
	}

	tmpl.Execute(w, nil)
}

// BlindTestHandler - Jeu Blind Test
// TODO @Quoc Huy: Charger les infos de la salle et du jeu
// TODO @Quoc Huy: Gérer la sélection de playlist (Rock/Rap/Pop)
// TODO @Quoc Huy: Timer 37s, système de points (3/2/1)
func BlindTestHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("./templates/games/blindtest.html", "./templates/header.html", "./templates/footer.html")
	if err != nil {
		log.Printf("Erreur: %v", err)
		http.Error(w, "Erreur serveur", http.StatusInternalServerError)
		return
	}
	// TODO: Passer les données du jeu au template
	tmpl.Execute(w, nil)
}

// PetitBacHandler - Jeu Petit Bac
// TODO @ilian: Charger les infos de la salle et du jeu
// TODO @ilian: Gérer les 9 manches, validation 2/3 joueurs
// TODO @ilian: Système de points (unique=2, commun=1)
func PetitBacHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("./templates/games/petitbac.html", "./templates/header.html", "./templates/footer.html")
	if err != nil {
		log.Printf("Erreur: %v", err)
		http.Error(w, "Erreur serveur", http.StatusInternalServerError)
		return
	}
	// TODO: Passer les données du jeu au template
	tmpl.Execute(w, nil)
}

var authService *services.AuthService
var roomService *services.RoomService

func main() {
	// TODO @Nome: Initialiser la connexion SQLite
	cfg := config.Load()

	db := database.InitDB(cfg.DatabasePath)
	defer db.Close()

	authService = &services.AuthService{DB: db}
	roomService = &services.RoomService{DB: db}

	// Initialiser les handlers avec les services
	handlers.Init(authService, roomService)

	// Routes
	http.HandleFunc("/", Landing)
	http.HandleFunc("/home", Home)
	http.HandleFunc("/register", RegisterHandler)
	http.HandleFunc("/login", LoginHandler)
	http.HandleFunc("/room/create", CreateRoomHandler)
	http.HandleFunc("/room/join", JoinRoomHandler)
	http.HandleFunc("/room/lobby", LobbyHandler)
	http.HandleFunc("/game/blindtest", BlindTestHandler)
	http.HandleFunc("/game/petitbac", PetitBacHandler)

	// TODO @Quoc Huy & @ilian: Routes API pour les actions de jeu
	// http.HandleFunc("/api/blindtest/submit-answer", SubmitBlindTestAnswerHandler)
	// http.HandleFunc("/api/petitbac/submit-answers", SubmitPetitBacAnswersHandler)
	// http.HandleFunc("/ws", WebSocketHandler)

	// Fichiers statiques
	fs := http.FileServer(http.Dir("static/"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	log.Println("🎵 Serveur Groupie Tracker sur http://localhost:" + cfg.ServerPort)
	log.Fatal(http.ListenAndServe(":"+cfg.ServerPort, nil))
}
