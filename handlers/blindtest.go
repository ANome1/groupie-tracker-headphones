package handlers

// RESPONSABLE: @Quoc Huy
// Handlers pour le jeu Blind Test

import "net/http"

// BlindTestGameHandler - Affiche l'interface du jeu Blind Test
// TODO @Quoc Huy:
//   - Vérifier que le joueur est dans la salle
//   - Afficher l'interface de jeu (templates/games/blindtest.html):
//     → Sélecteur de playlist: "Choisis une des playlist !"
//   - Rock
//   - Rap
//   - Pop
//     → Lecteur audio (lecture automatique de l'extrait)
//     → Timer visuel (37 secondes par défaut)
//     → Zone de saisie pour la réponse
//     → Scoreboard en temps réel
func BlindTestGameHandler(w http.ResponseWriter, r *http.Request) {
	// TODO: Implementation
}

// StartBlindTestHandler - Démarre une partie de Blind Test
// TODO @Quoc Huy:
// - Vérifier que l'utilisateur est l'hôte de la salle
// - Récupérer la playlist choisie (Rock, Rap, ou Pop)
// - Initialiser le jeu (models.NewBlindTestGame)
// - Récupérer les musiques depuis Spotify API
// - Changer le statut de la salle en 'in_progress'
// - Notifier tous les joueurs via WebSocket
// - Lancer le premier round
func StartBlindTestHandler(w http.ResponseWriter, r *http.Request) {
	// TODO: Implementation
}

// SubmitBlindTestAnswerHandler - Soumet une réponse à une musique
// TODO @Quoc Huy:
// - Récupérer la réponse du joueur
// - Récupérer le temps écoulé depuis le début du round
// - Vérifier la réponse (titre + artiste)
// - Calculer les points selon la rapidité:
//   - < 10s: 3 points
//   - < 20s: 2 points
//   - < 37s: 1 point
//
// - Mettre à jour le score
// - Notifier tous les joueurs via WebSocket
func SubmitBlindTestAnswerHandler(w http.ResponseWriter, r *http.Request) {
	// TODO: Implementation
}

// NextBlindTestRoundHandler - Passe au round suivant
// TODO @Quoc Huy:
// - Vérifier que l'utilisateur est l'hôte
// - Charger la musique suivante
// - Réinitialiser le timer (37 secondes)
// - Notifier tous les joueurs via WebSocket
// - Si tous les rounds sont terminés, afficher le scoreboard final
func NextBlindTestRoundHandler(w http.ResponseWriter, r *http.Request) {
	// TODO: Implementation
}
