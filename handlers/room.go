package handlers

import (
	"groupie-tracker/models"
	"html/template"
	"log"
	"net/http"
	"strconv"
)

type RoomData struct {
	Error        string
	Room         *models.Room
	Participants []models.RoomParticipant
	IsHost       bool
}

// CreateRoomHandler - Créer une nouvelle salle
func CreateRoomHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		tmpl, err := template.ParseFiles("./templates/room/create.html", "./templates/header.html", "./templates/footer.html")
		if err != nil {
			log.Printf("Erreur: %v", err)
			http.Error(w, "Erreur serveur", http.StatusInternalServerError)
			return
		}
		tmpl.Execute(w, RoomData{Error: ""})
		return
	}

	if r.Method == "POST" {
		// Vérifier que l'utilisateur est connecté
		cookie, err := r.Cookie("user_id")
		if err != nil {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		userID, err := strconv.Atoi(cookie.Value)
		if err != nil {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		// Récupérer les données du formulaire
		name := r.FormValue("name")
		gameType := r.FormValue("gameType") // "blindtest" ou "petitbac"

		if name == "" {
			tmpl, _ := template.ParseFiles("./templates/room/create.html", "./templates/header.html", "./templates/footer.html")
			tmpl.Execute(w, RoomData{Error: "Le nom de la salle est requis"})
			return
		}

		// Créer la salle
		room, err := RoomService.CreateRoom(name, gameType, userID)
		if err != nil {
			tmpl, _ := template.ParseFiles("./templates/room/create.html", "./templates/header.html", "./templates/footer.html")
			tmpl.Execute(w, RoomData{Error: "Erreur lors de la création: " + err.Error()})
			return
		}

		log.Printf("Salle créée: %s (code: %s, host: %d)", name, room.Code, userID)
		http.Redirect(w, r, "/room/lobby?code="+room.Code, http.StatusSeeOther)
	}
}

// JoinRoomHandler - Rejoindre une salle existante
func JoinRoomHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		tmpl, err := template.ParseFiles("./templates/room/join.html", "./templates/header.html", "./templates/footer.html")
		if err != nil {
			log.Printf("Erreur: %v", err)
			http.Error(w, "Erreur serveur", http.StatusInternalServerError)
			return
		}
		tmpl.Execute(w, RoomData{Error: ""})
		return
	}

	if r.Method == "POST" {
		// Vérifier que l'utilisateur est connecté
		cookie, err := r.Cookie("user_id")
		if err != nil {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		userID, err := strconv.Atoi(cookie.Value)
		if err != nil {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		// Récupérer le code de la salle
		roomCode := r.FormValue("roomCode")
		if roomCode == "" {
			tmpl, _ := template.ParseFiles("./templates/room/join.html", "./templates/header.html", "./templates/footer.html")
			tmpl.Execute(w, RoomData{Error: "Le code de la salle est requis"})
			return
		}

		// Trouver la salle
		room, err := RoomService.GetRoomByCode(roomCode)
		if err != nil {
			tmpl, _ := template.ParseFiles("./templates/room/join.html", "./templates/header.html", "./templates/footer.html")
			tmpl.Execute(w, RoomData{Error: "Salle non trouvée (code invalide)"})
			return
		}

		// Ajouter le participant
		err = RoomService.JoinRoom(room.ID, userID)
		if err != nil {
			tmpl, _ := template.ParseFiles("./templates/room/join.html", "./templates/header.html", "./templates/footer.html")
			tmpl.Execute(w, RoomData{Error: "Erreur lors de l'ajout: " + err.Error()})
			return
		}

		log.Printf("Utilisateur %d a rejoint la salle %s", userID, roomCode)
		http.Redirect(w, r, "/room/lobby?code="+roomCode, http.StatusSeeOther)
	}
}

// LobbyHandler - Afficher la salle et les participants
func LobbyHandler(w http.ResponseWriter, r *http.Request) {
	// Récupérer le code de la salle depuis les paramètres
	roomCode := r.URL.Query().Get("code")
	if roomCode == "" {
		http.Error(w, "Code de salle manquant", http.StatusBadRequest)
		return
	}

	// Récupérer la salle
	room, err := RoomService.GetRoomByCode(roomCode)
	if err != nil {
		http.Error(w, "Salle non trouvée", http.StatusNotFound)
		return
	}

	// Récupérer les participants
	participants, err := RoomService.GetRoomParticipants(room.ID)
	if err != nil {
		log.Printf("Erreur lors de la récupération des participants: %v", err)
		participants = []models.RoomParticipant{}
	}

	// Vérifier si l'utilisateur courant est le host
	cookie, _ := r.Cookie("user_id")
	userID := 0
	isHost := false
	if cookie != nil {
		userID, _ = strconv.Atoi(cookie.Value)
		isHost = (room.HostID == userID)
	}

	// Préparer les données
	data := RoomData{
		Room:         room,
		Participants: participants,
		IsHost:       isHost,
		Error:        "",
	}

	// Afficher le template
	tmpl, err := template.ParseFiles("./templates/room/lobby.html", "./templates/header.html", "./templates/footer.html")
	if err != nil {
		log.Printf("Erreur: %v", err)
		http.Error(w, "Erreur serveur", http.StatusInternalServerError)
		return
	}

	tmpl.Execute(w, data)
}

// LeaveRoomHandler - Quitter une salle
func LeaveRoomHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Méthode non autorisée", http.StatusMethodNotAllowed)
		return
	}

	// Vérifier que l'utilisateur est connecté
	cookie, err := r.Cookie("user_id")
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	userID, err := strconv.Atoi(cookie.Value)
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	// Récupérer le code de la salle
	roomCode := r.FormValue("roomCode")
	if roomCode == "" {
		http.Error(w, "Code de salle manquant", http.StatusBadRequest)
		return
	}

	// Trouver la salle
	room, err := RoomService.GetRoomByCode(roomCode)
	if err != nil {
		http.Error(w, "Salle non trouvée", http.StatusNotFound)
		return
	}

	// Supprimer le participant (à partir du roomID)
	// TODO: Implémenter RoomService.LeaveRoom(room.ID, userID)
	log.Printf("Utilisateur %d a quitté la salle %s (ID: %d)", userID, roomCode, room.ID)

	http.Redirect(w, r, "/home", http.StatusSeeOther)
}
