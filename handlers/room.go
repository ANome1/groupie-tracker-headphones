package handlers

import (
	"groupie-tracker/models"
	"groupie-tracker/utils"
	"html/template"
	"log"
	"net/http"
	"strconv"
)

type RoomData struct {
	User         *models.User
	Error        string
	Room         *models.Room
	Participants []models.RoomParticipant
	IsHost       bool
}

func CreateRoomHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		userID, _ := utils.GetUserIDFromCookie(r)
		var user *models.User
		if userID > 0 {
			user, _ = AuthService.GetUserByID(userID)
		}
		tmpl, err := template.ParseFiles("./templates/room/create.html", "./templates/header.html", "./templates/footer.html")
		if err != nil {
			log.Printf("Erreur: %v", err)
			http.Error(w, "Erreur serveur", http.StatusInternalServerError)
			return
		}
		tmpl.Execute(w, RoomData{User: user, Error: ""})
		return
	}

	if r.Method == "POST" {
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

		name := r.FormValue("name")
		gameType := r.FormValue("gameType") // "blindtest" ou "petitbac"

		if name == "" {
			tmpl, _ := template.ParseFiles("./templates/room/create.html", "./templates/header.html", "./templates/footer.html")
			tmpl.Execute(w, RoomData{Error: "Le nom de la salle est requis"})
			return
		}

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

func JoinRoomHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		tmpl, err := template.ParseFiles("./templates/room/join.html", "./templates/header.html", "./templates/footer.html")
		if err != nil {
			log.Printf("Erreur: %v", err)
			http.Error(w, "Erreur serveur", http.StatusInternalServerError)
			return
		}
		userID, _ := utils.GetUserIDFromCookie(r)
		var user *models.User
		if userID > 0 {
			user, _ = AuthService.GetUserByID(userID)
		}
		tmpl.Execute(w, RoomData{User: user, Error: ""})
		return
	}

	if r.Method == "POST" {
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

		roomCode := r.FormValue("roomCode")
		if roomCode == "" {
			tmpl, _ := template.ParseFiles("./templates/room/join.html", "./templates/header.html", "./templates/footer.html")
			tmpl.Execute(w, RoomData{Error: "Le code de la salle est requis"})
			return
		}

		room, err := RoomService.GetRoomByCode(roomCode)
		if err != nil {
			tmpl, _ := template.ParseFiles("./templates/room/join.html", "./templates/header.html", "./templates/footer.html")
			tmpl.Execute(w, RoomData{Error: "Salle non trouvée (code invalide)"})
			return
		}

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

func LobbyHandler(w http.ResponseWriter, r *http.Request) {
	roomCode := r.URL.Query().Get("code")
	if roomCode == "" {
		http.Error(w, "Code de salle manquant", http.StatusBadRequest)
		return
	}

	room, err := RoomService.GetRoomByCode(roomCode)
	if err != nil {
		http.Error(w, "Salle non trouvée", http.StatusNotFound)
		return
	}

	participants, err := RoomService.GetRoomParticipants(room.ID)
	if err != nil {
		log.Printf("Erreur lors de la récupération des participants: %v", err)
		participants = []models.RoomParticipant{}
	}

	cookie, _ := r.Cookie("user_id")
	userID := 0
	isHost := false
	if cookie != nil {
		userID, _ = strconv.Atoi(cookie.Value)
		isHost = (room.HostID == userID)
	}

	currentUserID, _ := utils.GetUserIDFromCookie(r)
	var user *models.User
	if currentUserID > 0 {
		user, _ = AuthService.GetUserByID(currentUserID)
	}
	data := RoomData{
		User:         user,
		Room:         room,
		Participants: participants,
		IsHost:       isHost,
		Error:        "",
	}

	tmpl, err := template.ParseFiles("./templates/room/lobby.html", "./templates/header.html", "./templates/footer.html")
	if err != nil {
		log.Printf("Erreur: %v", err)
		http.Error(w, "Erreur serveur", http.StatusInternalServerError)
		return
	}

	tmpl.Execute(w, data)
}

func LeaveRoomHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Méthode non autorisée", http.StatusMethodNotAllowed)
		return
	}

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

	roomCode := r.FormValue("roomCode")
	if roomCode == "" {
		http.Error(w, "Code de salle manquant", http.StatusBadRequest)
		return
	}

	room, err := RoomService.GetRoomByCode(roomCode)
	if err != nil {
		http.Error(w, "Salle non trouvée", http.StatusNotFound)
		return
	}

	err = RoomService.LeaveRoom(room.ID, userID)
	if err != nil {
		log.Printf("Erreur lors de la sortie de la salle: %v", err)
		http.Error(w, "Erreur serveur", http.StatusInternalServerError)
		return
	}

	log.Printf("Utilisateur %d a quitté la salle %s (ID: %d)", userID, roomCode, room.ID)

	http.Redirect(w, r, "/home", http.StatusSeeOther)
}
