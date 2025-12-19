package main

import (
	"groupie-tracker/config"
	"groupie-tracker/database"
	"groupie-tracker/handlers"
	"groupie-tracker/models"
	"groupie-tracker/services"
	"groupie-tracker/utils"
	"groupie-tracker/websocket"
	"io"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"text/template"
	"time"
)

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
func BlindTestHandler(w http.ResponseWriter, r *http.Request) {
	handlers.BlindTestHandler(w, r)
}

// PetitBacHandler - Jeu Petit Bac
func PetitBacHandler(w http.ResponseWriter, r *http.Request) {
	roomCode := r.URL.Query().Get("code")
	if roomCode == "" {
		http.Error(w, "Code de salle manquant", http.StatusBadRequest)
		return
	}

	room, err := roomService.GetRoomByCode(roomCode)
	if err != nil {
		http.Error(w, "Salle non trouvée", http.StatusNotFound)
		return
	}

	game := services.Manager.GetGame(roomCode)
	if game == nil {
		// Si le jeu n'existe pas encore (ex: refresh), on redirige vers le lobby ou on affiche une erreur
		// Pour l'instant, on redirige vers le lobby
		http.Redirect(w, r, "/room/lobby?code="+roomCode, http.StatusSeeOther)
		return
	}

	user := GetCurrentUser(r)

	var currentRound *models.PetitBacRound
	if len(game.Rounds) > 0 {
		currentRound = game.Rounds[len(game.Rounds)-1]
	}

	data := struct {
		User         *models.User
		Room         *models.Room
		Game         *models.PetitBacGame
		CurrentRound *models.PetitBacRound
	}{
		User:         user,
		Room:         room,
		Game:         game,
		CurrentRound: currentRound,
	}

	tmpl, err := template.ParseFiles("./templates/games/petitbac.html", "./templates/components/header.html", "./templates/components/footer.html")
	if err != nil {
		log.Printf("Erreur: %v", err)
		http.Error(w, "Erreur serveur", http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, data)
}

// StartGameHandler - Démarre la partie
func StartGameHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Méthode non autorisée", http.StatusMethodNotAllowed)
		return
	}

	cookie, err := r.Cookie("user_id")
	if err != nil {
		http.Error(w, "Non authentifié", http.StatusUnauthorized)
		return
	}

	userID, _ := strconv.Atoi(cookie.Value)
	roomCode := r.FormValue("roomCode")

	// Parse configuration
	numRounds, _ := strconv.Atoi(r.FormValue("numRounds"))
	if numRounds <= 0 {
		numRounds = 5 // Default
	}

	timePerRound, _ := strconv.Atoi(r.FormValue("timePerRound"))
	if timePerRound <= 0 {
		timePerRound = 30 // Default
	}

	room, err := roomService.GetRoomByCode(roomCode)
	if err != nil {
		http.Error(w, "Salle non trouvée", http.StatusNotFound)
		return
	}

	if room.HostID != userID {
		http.Error(w, "Seul l'hôte peut démarrer la partie", http.StatusForbidden)
		return
	}

	// Initialiser le jeu selon le type
	if room.GameType == "petitbac" {
		// Force timePerRound to 60 for Petit Bac
		timePerRound = 60

		// Récupérer les participants
		participants, err := roomService.GetRoomParticipantsWithUsers(room.ID)
		if err != nil {
			http.Error(w, "Erreur lors de la récupération des participants", http.StatusInternalServerError)
			return
		}

		players := make([]string, len(participants))
		playerNames := make(map[string]string)
		for i, p := range participants {
			pid := strconv.Itoa(p.UserID)
			players[i] = pid
			if p.Username != "" {
				playerNames[pid] = p.Username
			} else {
				playerNames[pid] = "Joueur " + pid
			}
		}

		// Config Petit Bac
		config := models.PetitBacConfig{
			Categories:   []string{"Artiste", "Groupe de musique", "Album", "Instrument", "Featuring"},
			TimePerRound: timePerRound,
			NumRounds:    numRounds,
		}

		// Créer l'instance de jeu
		game := services.Manager.CreateGame(room.Code, strconv.Itoa(room.HostID), players, playerNames, config)

		// Démarrer le premier round immédiatement pour avoir une lettre
		services.Manager.StartRound(game.ID)

		// Diffuser le message de début de partie via WebSocket
		hub.BroadcastToRoom(room.Code, models.MessageOut{
			Type: "GAME_START",
			Data: "/game/petitbac?code=" + room.Code,
		})
	} else if room.GameType == "blindtest" {
		playlistID := r.FormValue("playlistId")

		// Créer une partie de blind test AVEC la configuration
		config := models.BlindTestConfig{
			PlaylistID:   playlistID,
			TimePerRound: timePerRound,
			NumRounds:    numRounds,
		}
		services.BlindTestMgr.CreateGame(room.Code, room.HostID, config)

		// Démarrer la première manche immédiatement SEULEMENT si une playlist a été sélectionnée
		if playlistID != "" {
			game := services.BlindTestMgr.GetGame(room.Code)
			if game != nil {
				// Utiliser une goroutine pour ne pas bloquer la réponse HTTP
				// et laisser le temps aux clients de se connecter au WebSocket
				go func() {
					time.Sleep(3 * time.Second) // Petit délai pour la connexion WS
					services.BlindTestMgr.StartGameRound(game, hub.BroadcastToRoom)
				}()
			}
		}

		// Diffuser le message de début de partie via WebSocket
		hub.BroadcastToRoom(room.Code, models.MessageOut{
			Type: "GAME_START",
			Data: "/game/blindtest?code=" + room.Code,
		})
		log.Printf("Blind test initialisé pour room %s avec config: %+v", room.Code, config)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}

// AudioProxyHandler proxies audio files from Deezer to avoid CORS issues
func AudioProxyHandler(w http.ResponseWriter, r *http.Request) {
	// Get the audio URL from query parameter
	audioURL := r.URL.Query().Get("url")
	if audioURL == "" {
		http.Error(w, "Missing url parameter", http.StatusBadRequest)
		return
	}

	// Validate that the URL is from Deezer (security check)
	parsedURL, err := url.Parse(audioURL)
	if err != nil {
		http.Error(w, "Invalid URL", http.StatusBadRequest)
		return
	}

	// Allow both cdnt-preview.dzcdn.net and cdns-files-c.dzcdn.net
	if parsedURL.Host != "cdnt-preview.dzcdn.net" && parsedURL.Host != "cdns-files-c.dzcdn.net" {
		http.Error(w, "Invalid or untrusted URL", http.StatusBadRequest)
		return
	}

	// Fetch the audio from Deezer
	resp, err := http.Get(audioURL)
	if err != nil {
		log.Printf("Error fetching audio from Deezer: %v", err)
		http.Error(w, "Failed to fetch audio", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	// Check response status
	if resp.StatusCode != http.StatusOK {
		http.Error(w, "Failed to fetch audio", http.StatusInternalServerError)
		return
	}

	// Set proper headers for audio streaming
	w.Header().Set("Content-Type", resp.Header.Get("Content-Type"))
	if resp.Header.Get("Content-Length") != "" {
		w.Header().Set("Content-Length", resp.Header.Get("Content-Length"))
	}
	w.Header().Set("Accept-Ranges", "bytes")
	w.Header().Set("Cache-Control", "public, max-age=3600")

	// Stream the audio
	_, err = io.Copy(w, resp.Body)
	if err != nil {
		log.Printf("Error streaming audio: %v", err)
	}
}

var authService *services.AuthService
var roomService *services.RoomService
var hub *websocket.Hub

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
	cfg := config.Load()

	db := database.InitDB(cfg.DatabasePath)
	defer db.Close()

	authService = &services.AuthService{DB: db}
	roomService = &services.RoomService{DB: db}

	// Initialiser le Hub WebSocket
	hub = websocket.NewHub()
	go hub.Run()

	// Initialiser les handlers avec les services et le hub
	handlers.Init(authService, roomService)
	handlers.SetGlobalHub(hub)

	// Middleware de sécurité pour tous les handlers
	securityHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// En-têtes de sécurité
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("X-Frame-Options", "DENY")
		w.Header().Set("X-XSS-Protection", "1; mode=block")
		w.Header().Set("Referrer-Policy", "strict-origin-when-cross-origin")
		
		// HSTS (Strict-Transport-Security) - activer en production avec HTTPS
		// w.Header().Set("Strict-Transport-Security", "max-age=31536000; includeSubDomains; preload")
	})

	// Routes
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		securityHandler.ServeHTTP(w, r)
		Landing(w, r)
	})
	http.HandleFunc("/home", func(w http.ResponseWriter, r *http.Request) {
		securityHandler.ServeHTTP(w, r)
		Home(w, r)
	})
	http.HandleFunc("/register", func(w http.ResponseWriter, r *http.Request) {
		securityHandler.ServeHTTP(w, r)
		RegisterHandler(w, r)
	})
	http.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		securityHandler.ServeHTTP(w, r)
		LoginHandler(w, r)
	})
	http.HandleFunc("/logout", func(w http.ResponseWriter, r *http.Request) {
		securityHandler.ServeHTTP(w, r)
		LogoutHandler(w, r)
	})
	http.HandleFunc("/room/create", func(w http.ResponseWriter, r *http.Request) {
		securityHandler.ServeHTTP(w, r)
		CreateRoomHandler(w, r)
	})
	http.HandleFunc("/room/join", func(w http.ResponseWriter, r *http.Request) {
		securityHandler.ServeHTTP(w, r)
		JoinRoomHandler(w, r)
	})
	http.HandleFunc("/room/lobby", func(w http.ResponseWriter, r *http.Request) {
		securityHandler.ServeHTTP(w, r)
		LobbyHandler(w, r)
	})
	http.HandleFunc("/room/leave", func(w http.ResponseWriter, r *http.Request) {
		securityHandler.ServeHTTP(w, r)
		LeaveRoomHandler(w, r)
	})
	http.HandleFunc("/game/blindtest", func(w http.ResponseWriter, r *http.Request) {
		securityHandler.ServeHTTP(w, r)
		BlindTestHandler(w, r)
	})
	http.HandleFunc("/game/petitbac", func(w http.ResponseWriter, r *http.Request) {
		securityHandler.ServeHTTP(w, r)
		PetitBacHandler(w, r)
	})
	http.HandleFunc("/game/start", func(w http.ResponseWriter, r *http.Request) {
		securityHandler.ServeHTTP(w, r)
		StartGameHandler(w, r)
	})
	http.HandleFunc("/audio/proxy", func(w http.ResponseWriter, r *http.Request) {
		securityHandler.ServeHTTP(w, r)
		AudioProxyHandler(w, r)
	})

	// WebSocket
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		websocket.ServeWs(hub, w, r)
	})

	// Fichiers statiques
	fs := http.FileServer(http.Dir("static/"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	log.Println("🎵 Serveur Groupie Tracker sur http://localhost:" + cfg.ServerPort)
	log.Fatal(http.ListenAndServe(":"+cfg.ServerPort, nil))
}
