package handlers

import "net/http"

// RESPONSABLE: @Nome (avec collaboration de @Quoc Huy et @ilian)
// Handlers pour la gestion des salles de jeu

// TODO @Nome: CreateRoomHandler - Créer une nouvelle salle
// - Méthode GET: Afficher le formulaire de création (templates/room/create.html)
// - Méthode POST:
//   1. Récupérer le nom de la salle, type de jeu (blindtest/petitbac), max joueurs
//   2. Vérifier que l'utilisateur est connecté (session)
//   3. Générer un code de salle unique (6 caractères)
//   4. Créer la salle en base de données
//   5. Ajouter le créateur comme premier participant
//   6. Rediriger vers le lobby de la salle

// TODO @Nome: JoinRoomHandler - Rejoindre une salle existante
// - Méthode GET: Afficher le formulaire pour entrer le code (templates/room/join.html)
// - Méthode POST:
//   1. Récupérer le code de la salle
//   2. Vérifier que l'utilisateur est connecté
//   3. Vérifier que la salle existe et n'est pas pleine
//   4. Ajouter l'utilisateur comme participant
//   5. Rediriger vers le lobby de la salle

// TODO @Nome: LobbyHandler - Affiche le lobby d'attente avant le jeu
// - Afficher la liste des participants
// - Afficher le code de la salle pour partager
// - Bouton "Commencer" visible uniquement pour l'hôte
// - WebSocket pour mettre à jour la liste en temps réel

// TODO @Nome: LeaveRoomHandler - Quitter une salle
// - Retirer l'utilisateur des participants
// - Si c'était l'hôte, transférer à un autre joueur ou supprimer la salle
// - Notifier les autres joueurs via WebSocket

// LeaveRoomHandler - Quitter une salle
// TODO @Nome:
// - Retirer le joueur de la salle
// - Si c'était l'hôte, désigner un nouvel hôte ou fermer la salle
// - Notifier les autres joueurs via WebSocket
func LeaveRoomHandler(w http.ResponseWriter, r *http.Request) {
	// TODO: Implementation
}
