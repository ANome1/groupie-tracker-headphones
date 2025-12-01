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
		// TODO: Logique d'inscription
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	tmpl, err := template.ParseFiles("./templates/auth/register.html")
	if err != nil {
		log.Printf("Erreur: %v", err)
		http.Error(w, "Erreur serveur", http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, nil)
}

// LoginHandler - Connexion
// TODO @Nome: Récupérer username, password du formulaire
// TODO @Nome: Vérifier en base, comparer les hash
// TODO @Nome: Créer une session si succès
func LoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		// TODO: Logique de connexion
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	tmpl, err := template.ParseFiles("./templates/auth/login.html")
	if err != nil {
		log.Printf("Erreur: %v", err)
		http.Error(w, "Erreur serveur", http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, nil)
}

// CreateRoomHandler - Créer une salle
// TODO @Nome: Récupérer room_name, game_type (blindtest/petitbac), max_players
// TODO @Nome: Générer un code unique (6 caractères)
// TODO @Nome: Insérer en base, ajouter le créateur comme participant
func CreateRoomHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		// TODO: Logique création salle
		http.Redirect(w, r, "/room/lobby", http.StatusSeeOther)
		return
	}

	tmpl, err := template.ParseFiles("./templates/room/create.html")
	if err != nil {
		log.Printf("Erreur: %v", err)
		http.Error(w, "Erreur serveur", http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, nil)
}

// JoinRoomHandler - Rejoindre une salle
// TODO @Nome: Récupérer room_code du formulaire
// TODO @Nome: Vérifier que la salle existe et n'est pas pleine
// TODO @Nome: Ajouter le joueur aux participants
func JoinRoomHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		// TODO: Logique rejoindre salle
		http.Redirect(w, r, "/room/lobby", http.StatusSeeOther)
		return
	}

	tmpl, err := template.ParseFiles("./templates/room/join.html")
	if err != nil {
		log.Printf("Erreur: %v", err)
		http.Error(w, "Erreur serveur", http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, nil)
}

// LobbyHandler - Lobby d'attente
// TODO @Nome: Charger les participants de la salle depuis la base
// TODO @Nome: Afficher le code de salle, liste des joueurs
// TODO @Nome: Bouton "Commencer" uniquement pour l'hôte
func LobbyHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("./templates/room/lobby.html")
	if err != nil {
		log.Printf("Erreur: %v", err)
		http.Error(w, "Erreur serveur", http.StatusInternalServerError)
		return
	}
	// TODO: Passer les données de la salle au template
	tmpl.Execute(w, nil)
}

// BlindTestHandler - Jeu Blind Test
// TODO @Quoc Huy: Charger les infos de la salle et du jeu
// TODO @Quoc Huy: Gérer la sélection de playlist (Rock/Rap/Pop)
// TODO @Quoc Huy: Timer 37s, système de points (3/2/1)
func BlindTestHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("./templates/games/blindtest.html")
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
	tmpl, err := template.ParseFiles("./templates/games/petitbac.html")
	if err != nil {
		log.Printf("Erreur: %v", err)
		http.Error(w, "Erreur serveur", http.StatusInternalServerError)
		return
	}
	// TODO: Passer les données du jeu au template
	tmpl.Execute(w, nil)
}

func main() {
	// TODO @Nome: Initialiser la connexion SQLite
	// db, err := sql.Open("sqlite3", "./database/groupie-tracker.db")

	// Routes
	http.HandleFunc("/", Home)
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

	log.Println("🎵 Serveur Groupie Tracker sur http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}
