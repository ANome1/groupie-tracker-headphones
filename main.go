package main

import (
	"fmt"
	"log"
	"net/http"
)

// RESPONSABLE: @Nome (infrastructure), @Quoc Huy (WebSocket Blind Test), @ilian (WebSocket Petit Bac)
// Ce fichier initialise le serveur web, les routes et les connexions

func main() {
	// TODO @Nome: Initialiser la connexion à la base de données SQLite
	// - Charger le fichier database/groupie-tracker.db
	// - Exécuter le schéma SQL si la base n'existe pas

	// Servir les fichiers statiques (CSS, JS, images)
	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	// TODO @Nome: Routes d'authentification
	http.HandleFunc("/register", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "templates/auth/register.html")
	})
	http.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "templates/auth/login.html")
	})

	// TODO @Nome: Routes de gestion des salles
	http.HandleFunc("/room/create", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "templates/room/create.html")
	})
	http.HandleFunc("/room/join", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "templates/room/join.html")
	})
	http.HandleFunc("/room/lobby", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "templates/room/lobby.html")
	})

	// TODO @Quoc Huy: Routes du Blind Test
	http.HandleFunc("/game/blindtest", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "templates/games/blindtest.html")
	})

	// TODO @ilian: Routes du Petit Bac
	http.HandleFunc("/game/petitbac", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "templates/games/petitbac.html")
	})

	// Landing page
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}
		http.ServeFile(w, r, "templates/landing.html")
	})

	// TODO @Quoc Huy & @ilian: Initialiser le hub WebSocket
	// - Lancer le hub.Run() en goroutine
	// - Route /ws pour les connexions WebSocket

	fmt.Println("🎵 Serveur Groupie Tracker démarré sur http://localhost:8080")
	fmt.Println("Appuyez sur Ctrl+C pour arrêter le serveur")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
