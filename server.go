package main

import (
	"log"
	"net/http"
	"text/template"
)

// RESPONSABLE: @Nome (infrastructure), @Quoc Huy (WebSocket Blind Test), @ilian (WebSocket Petit Bac)

// Home - Page d'accueil
func Home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	tmpl, err := template.ParseFiles("./templates/index.html")
	if err != nil {
		log.Printf("Erreur lors du chargement du template: %v", err)
		http.Error(w, "Erreur serveur", http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, nil)
}

// RegisterHandler - Page d'inscription
// TODO @Nome: Implémenter la logique d'inscription
func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		// TODO: Traiter le formulaire d'inscription
		// username := r.FormValue("username")
		// email := r.FormValue("email")
		// password := r.FormValue("password")
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	tmpl, err := template.ParseFiles("./templates/auth/register.html")
	if err != nil {
		log.Printf("Erreur template: %v", err)
		http.Error(w, "Erreur serveur", http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, nil)
}

// LoginHandler - Page de connexion
// TODO @Nome: Implémenter la logique de connexion
func LoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		// TODO: Traiter le formulaire de connexion
		// username := r.FormValue("username")
		// password := r.FormValue("password")
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	tmpl, err := template.ParseFiles("./templates/auth/login.html")
	if err != nil {
		log.Printf("Erreur template: %v", err)
		http.Error(w, "Erreur serveur", http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, nil)
}

// CreateRoomHandler - Créer une salle
// TODO @Nome: Implémenter la logique de création de salle
func CreateRoomHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		// TODO: Créer la salle
		// roomName := r.FormValue("room_name")
		// gameType := r.FormValue("game_type")
		// maxPlayers := r.FormValue("max_players")
		http.Redirect(w, r, "/room/lobby", http.StatusSeeOther)
		return
	}

	tmpl, err := template.ParseFiles("./templates/room/create.html")
	if err != nil {
		log.Printf("Erreur template: %v", err)
		http.Error(w, "Erreur serveur", http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, nil)
}

// JoinRoomHandler - Rejoindre une salle
// TODO @Nome: Implémenter la logique pour rejoindre une salle
func JoinRoomHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		// TODO: Rejoindre la salle
		// roomCode := r.FormValue("room_code")
		http.Redirect(w, r, "/room/lobby", http.StatusSeeOther)
		return
	}

	tmpl, err := template.ParseFiles("./templates/room/join.html")
	if err != nil {
		log.Printf("Erreur template: %v", err)
		http.Error(w, "Erreur serveur", http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, nil)
}

// LobbyHandler - Lobby d'attente
// TODO @Nome: Implémenter la logique du lobby
func LobbyHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("./templates/room/lobby.html")
	if err != nil {
		log.Printf("Erreur template: %v", err)
		http.Error(w, "Erreur serveur", http.StatusInternalServerError)
		return
	}

	// TODO: Charger les données de la salle
	// data := struct {
	// 	RoomCode string
	// 	Players  []string
	// 	IsHost   bool
	// }{}

	tmpl.Execute(w, nil)
}

// BlindTestHandler - Page du jeu Blind Test
// TODO @Quoc Huy: Implémenter la logique du Blind Test
func BlindTestHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("./templates/games/blindtest.html")
	if err != nil {
		log.Printf("Erreur template: %v", err)
		http.Error(w, "Erreur serveur", http.StatusInternalServerError)
		return
	}

	// TODO: Charger les données du jeu
	tmpl.Execute(w, nil)
}

// PetitBacHandler - Page du jeu Petit Bac
// TODO @ilian: Implémenter la logique du Petit Bac
func PetitBacHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("./templates/games/petitbac.html")
	if err != nil {
		log.Printf("Erreur template: %v", err)
		http.Error(w, "Erreur serveur", http.StatusInternalServerError)
		return
	}

	// TODO: Charger les données du jeu
	tmpl.Execute(w, nil)
}

func main() {
	// TODO @Nome: Initialiser la connexion à la base de données SQLite
	// - Charger le fichier database/groupie-tracker.db
	// - Exécuter le schéma SQL si la base n'existe pas

	// Routes principales
	http.HandleFunc("/", Home)

	// Routes d'authentification
	http.HandleFunc("/register", RegisterHandler)
	http.HandleFunc("/login", LoginHandler)

	// Routes des salles
	http.HandleFunc("/room/create", CreateRoomHandler)
	http.HandleFunc("/room/join", JoinRoomHandler)
	http.HandleFunc("/room/lobby", LobbyHandler)

	// Routes des jeux
	http.HandleFunc("/game/blindtest", BlindTestHandler)
	http.HandleFunc("/game/petitbac", PetitBacHandler)

	// TODO @Quoc Huy & @ilian: Routes API et WebSocket
	// http.HandleFunc("/api/blindtest/submit", SubmitAnswerHandler)
	// http.HandleFunc("/api/petitbac/submit", SubmitAnswersHandler)
	// http.HandleFunc("/ws", WebSocketHandler)

	// Servir les fichiers statiques (CSS, JS, images)
	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	log.Println("🎵 Serveur Groupie Tracker démarré sur http://localhost:8080")
	log.Println("Appuyez sur Ctrl+C pour arrêter le serveur")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
