package handlers

import (
	"encoding/json"
	"fmt"
	"groupie-tracker/models"
	"groupie-tracker/utils"
	"groupie-tracker/websocket"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"strings"
)

var globalHub *websocket.Hub

func SetGlobalHub(hub *websocket.Hub) {
	globalHub = hub
}

type RoomData struct {
	User         *models.User
	Error        string
	Room         *models.Room
	Participants []models.RoomParticipantWithUser
	IsHost       bool
}

func CreateRoomHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		userID, _ := utils.GetUserIDFromCookie(r)
		var user *models.User
		if userID > 0 {
			user, _ = AuthService.GetUserByID(userID)
		}
		tmpl, err := template.ParseFiles("./templates/room/create.html", "./templates/components/header.html", "./templates/components/footer.html")
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
		gameType := r.FormValue("gameType")
		maxPlayersStr := r.FormValue("max_players")

		if name == "" {
			tmpl, _ := template.ParseFiles("./templates/room/create.html", "./templates/components/header.html", "./templates/components/footer.html")
			tmpl.Execute(w, RoomData{Error: "Le nom de la salle est requis"})
			return
		}

		maxPlayers := 8 // Valeur par défaut
		if maxPlayersStr != "" {
			if parsedMax, err := strconv.Atoi(maxPlayersStr); err == nil {
				maxPlayers = parsedMax
			}
		}

		room, err := RoomService.CreateRoom(name, gameType, userID, maxPlayers)
		if err != nil {
			tmpl, _ := template.ParseFiles("./templates/room/create.html", "./templates/components/header.html", "./templates/components/footer.html")
			tmpl.Execute(w, RoomData{Error: "Erreur lors de la création: " + err.Error()})
			return
		}

		err = RoomService.JoinRoom(room.ID, userID)
		if err != nil {
			log.Printf("Erreur lors de l'ajout de l'hôte: %v", err)
		}

		log.Printf("Salle créée: %s (code: %s, host: %d, max_players: %d)", name, room.Code, userID, maxPlayers)
		http.Redirect(w, r, "/room/lobby?code="+room.Code, http.StatusSeeOther)
	}
}

func JoinRoomHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		tmpl, err := template.ParseFiles("./templates/room/join.html", "./templates/components/header.html", "./templates/components/footer.html")
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
			tmpl, _ := template.ParseFiles("./templates/room/join.html", "./templates/components/header.html", "./templates/components/footer.html")
			tmpl.Execute(w, RoomData{Error: "Le code de la salle est requis"})
			return
		}

		roomCode = strings.ToUpper(strings.TrimSpace(roomCode))

		log.Printf("Tentative de connexion à la salle: '%s' par l'utilisateur %d", roomCode, userID)

		room, err := RoomService.GetRoomByCode(roomCode)
		if err != nil {
			log.Printf("Erreur GetRoomByCode: %v", err)
			tmpl, _ := template.ParseFiles("./templates/room/join.html", "./templates/components/header.html", "./templates/components/footer.html")
			tmpl.Execute(w, RoomData{Error: "Salle non trouvée (code invalide)"})
			return
		}

		// Vérifier si la salle est pleine
		participants, err := RoomService.GetRoomParticipantsWithUsers(room.ID)
		if err != nil {
			log.Printf("Erreur récupération participants: %v", err)
			tmpl, _ := template.ParseFiles("./templates/room/join.html", "./templates/components/header.html", "./templates/components/footer.html")
			tmpl.Execute(w, RoomData{Error: "Erreur serveur lors de la vérification de la salle"})
			return
		}

		if len(participants) >= room.MaxPlayers {
			tmpl, _ := template.ParseFiles("./templates/room/join.html", "./templates/components/header.html", "./templates/components/footer.html")
			tmpl.Execute(w, RoomData{Error: fmt.Sprintf("La salle est pleine (%d/%d joueurs)", len(participants), room.MaxPlayers)})
			return
		}

		err = RoomService.JoinRoom(room.ID, userID)
		if err != nil {
			log.Printf("Erreur JoinRoom: %v", err)
			tmpl, _ := template.ParseFiles("./templates/room/join.html", "./templates/components/header.html", "./templates/components/footer.html")
			tmpl.Execute(w, RoomData{Error: "Erreur lors de l'ajout: " + err.Error()})
			return
		}

		log.Printf("Utilisateur %d a rejoint la salle %s", userID, roomCode)

		user, err := AuthService.GetUserByID(userID)
		userName := "Utilisateur"
		if err == nil && user != nil {
			userName = user.Username
		}

		// Récupérer les participants mis à jour après la jointure
		participants, err = RoomService.GetRoomParticipantsWithUsers(room.ID)
		if err != nil {
			log.Printf("Erreur récupération participants après jointure: %v", err)
			participants = []models.RoomParticipantWithUser{}
		}

		if globalHub != nil {
			globalHub.BroadcastToRoom(roomCode, models.MessageOut{
				Type: "PLAYER_JOINED",
				Data: map[string]interface{}{
					"username":     userName,
					"userID":       userID,
					"participants": participants,
					"total":        len(participants),
				},
			})
		}

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

	participants, err := RoomService.GetRoomParticipantsWithUsers(room.ID)
	if err != nil {
		log.Printf("Erreur lors de la récupération des participants: %v", err)
		participants = []models.RoomParticipantWithUser{}
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

	tmpl, err := template.ParseFiles("./templates/room/lobby.html", "./templates/components/header.html", "./templates/components/footer.html")
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

	user, _ := AuthService.GetUserByID(userID)
	userName := "Utilisateur"
	if user != nil {
		userName = user.Username
	}

	participants, _ := RoomService.GetRoomParticipantsWithUsers(room.ID)

	if globalHub != nil {
		globalHub.BroadcastToRoom(roomCode, models.MessageOut{
			Type: "PLAYER_LEFT",
			Data: map[string]interface{}{
				"username":     userName,
				"userID":       userID,
				"participants": participants,
				"total":        len(participants),
			},
		})
	}

	http.Redirect(w, r, "/home", http.StatusSeeOther)
}

func ChangeGameTypeHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Get user ID from cookie
	cookie, err := r.Cookie("user_id")
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	userID, err := strconv.Atoi(cookie.Value)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	// Parse JSON body
	var req struct {
		RoomCode string `json:"roomCode"`
		GameType string `json:"gameType"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	// Get the room
	room, err := RoomService.GetRoomByCode(req.RoomCode)
	if err != nil {
		http.Error(w, "Room not found", http.StatusNotFound)
		return
	}

	// Verify the user is the host
	if room.HostID != userID {
		http.Error(w, "Only the host can change the game type", http.StatusForbidden)
		return
	}

	// Validate game type
	if req.GameType != "blindtest" && req.GameType != "petitbac" {
		http.Error(w, "Invalid game type", http.StatusBadRequest)
		return
	}

	// Update the game type
	err = RoomService.UpdateRoomGameType(room.ID, req.GameType)
	if err != nil {
		http.Error(w, "Error updating game type: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Notify all players in the room via WebSocket
	globalHub.BroadcastToRoom(req.RoomCode, map[string]interface{}{
		"Type": "GAME_TYPE_CHANGED",
		"Data": map[string]interface{}{
			"gameType": req.GameType,
		},
	})

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, `{"message":"Game type updated"}`)
}
