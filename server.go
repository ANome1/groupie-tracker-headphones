package main

import (
	"groupie-tracker/config"
	"groupie-tracker/database"
	"groupie-tracker/handlers"
	"groupie-tracker/models"
	"groupie-tracker/services"
	"groupie-tracker/utils"
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

	user := GetCurrentUser(r)
	if user != nil {
		http.Redirect(w, r, "/home", http.StatusSeeOther)
		return
	}

	data := struct {
		User *models.User
	}{User: user}

	tmpl, err := template.ParseFiles("./templates/landing.html", "./templates/components/header.html", "./templates/components/footer.html")
	if err != nil {
		log.Printf("Erreur: %v", err)
		http.Error(w, "Erreur serveur", http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, data)
}

// Home - Page de sélection de jeu
func Home(w http.ResponseWriter, r *http.Request) {
	user := GetCurrentUser(r)
	data := struct {
		User *models.User
	}{User: user}

	tmpl, err := template.ParseFiles("./templates/home.html", "./templates/components/header.html", "./templates/components/footer.html")
	if err != nil {
		log.Printf("Erreur: %v", err)
		http.Error(w, "Erreur serveur", http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, data)
}

// RegisterHandler - Inscription
func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	handlers.RegisterHandler(w, r)
}

// LoginHandler - Connexion
func LoginHandler(w http.ResponseWriter, r *http.Request) {
	handlers.LoginHandler(w, r)
}

// LogoutHandler - Déconnexion
func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	handlers.LogoutHandler(w, r)
}

// CreateRoomHandler - Créer une salle
func CreateRoomHandler(w http.ResponseWriter, r *http.Request) {
	handlers.CreateRoomHandler(w, r)
}

// JoinRoomHandler - Rejoindre une salle
func JoinRoomHandler(w http.ResponseWriter, r *http.Request) {
	handlers.JoinRoomHandler(w, r)
}

// LobbyHandler - Afficher la salle
func LobbyHandler(w http.ResponseWriter, r *http.Request) {
	handlers.LobbyHandler(w, r)
}

// LeaveRoomHandler - Quitter une salle
func LeaveRoomHandler(w http.ResponseWriter, r *http.Request) {
	handlers.LeaveRoomHandler(w, r)
}

// BlindTestHandler - Jeu Blind Test
// TODO @Quoc Huy: Charger les infos de la salle et du jeu
// TODO @Quoc Huy: Gérer la sélection de playlist (Rock/Rap/Pop)
// TODO @Quoc Huy: Timer 37s, système de points (3/2/1)
func BlindTestHandler(w http.ResponseWriter, r *http.Request) {
	user := GetCurrentUser(r)
	data := struct {
		User *models.User
	}{User: user}

	tmpl, err := template.ParseFiles("./templates/games/blindtest.html", "./templates/components/header.html", "./templates/components/footer.html", "./templates/components/scoreboard.html")
	if err != nil {
		log.Printf("Erreur: %v", err)
		http.Error(w, "Erreur serveur", http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, data)
}

// PetitBacHandler - Jeu Petit Bac
// TODO @ilian: Charger les infos de la salle et du jeu
// TODO @ilian: Gérer les 9 manches, validation 2/3 joueurs
// TODO @ilian: Système de points (unique=2, commun=1)
func PetitBacHandler(w http.ResponseWriter, r *http.Request) {
	user := GetCurrentUser(r)
	data := struct {
		User *models.User
	}{User: user}

	tmpl, err := template.ParseFiles("./templates/games/petitbac.html", "./templates/components/header.html", "./templates/components/footer.html", "./templates/components/scoreboard.html")
	if err != nil {
		log.Printf("Erreur: %v", err)
		http.Error(w, "Erreur serveur", http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, data)
}

var authService *services.AuthService
var roomService *services.RoomService

func GetCurrentUser(r *http.Request) *models.User {
	userID, err := utils.GetUserIDFromCookie(r)
	if err != nil {
		return nil
	}

	user, err := authService.GetUserByID(userID)
	if err != nil {
		return nil
	}

	return user
}

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
	http.HandleFunc("/logout", LogoutHandler)
	http.HandleFunc("/room/create", CreateRoomHandler)
	http.HandleFunc("/room/join", JoinRoomHandler)
	http.HandleFunc("/room/lobby", LobbyHandler)
	http.HandleFunc("/room/leave", LeaveRoomHandler)
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
