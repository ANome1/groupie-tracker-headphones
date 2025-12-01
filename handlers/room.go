package handlers

// RESPONSABLE: @Nome (avec collaboration de @Quoc Huy et @ilian)
// Handlers pour la gestion des salles de jeu

import "net/http"

// CreateRoomHandler - Crée une nouvelle salle de jeu
// TODO @Nome:
//   - Méthode GET: Afficher le formulaire de création (templates/room/create.html)
//     → Avec options: Blind Test ou Petit Bac
//   - Méthode POST:
//     1. Récupérer: room_name, game_type, max_players
//     2. Générer un code de salle unique (6 caractères)
//     3. Créer la salle en base de données
//     4. Ajouter le créateur comme participant
//     5. Rediriger vers le lobby de la salle
func CreateRoomHandler(w http.ResponseWriter, r *http.Request) {
	// TODO: Implementation
}

// JoinRoomHandler - Permet de rejoindre une salle existante
// TODO @Nome:
//   - Méthode GET: Afficher le formulaire (templates/room/join.html)
//     → Champ pour entrer le code de la salle
//   - Méthode POST:
//     1. Récupérer le code de salle
//     2. Vérifier que la salle existe
//     3. Vérifier que la salle n'est pas pleine
//     4. Vérifier que la salle est en statut 'waiting'
//     5. Ajouter le joueur à la salle
//     6. Rediriger vers le lobby de la salle
func JoinRoomHandler(w http.ResponseWriter, r *http.Request) {
	// TODO: Implementation
}

// LobbyHandler - Affiche le lobby d'une salle
// TODO @Nome:
//   - Afficher les informations de la salle (templates/room/lobby.html):
//     → Nom de la salle
//     → Code de la salle (pour inviter d'autres joueurs)
//     → Type de jeu
//     → Liste des participants
//     → Bouton "Démarrer" (seulement pour l'hôte)
//   - Utiliser WebSocket pour mettre à jour la liste des joueurs en temps réel
func LobbyHandler(w http.ResponseWriter, r *http.Request) {
	// TODO: Implementation
}

// LeaveRoomHandler - Quitter une salle
// TODO @Nome:
// - Retirer le joueur de la salle
// - Si c'était l'hôte, désigner un nouvel hôte ou fermer la salle
// - Notifier les autres joueurs via WebSocket
func LeaveRoomHandler(w http.ResponseWriter, r *http.Request) {
	// TODO: Implementation
}
