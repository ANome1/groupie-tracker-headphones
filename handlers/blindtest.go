package handlers

// RESPONSABLE: @Quoc Huy
// Handlers pour le jeu Blind Test

// TODO @Quoc Huy: BlindTestHandler - Page principale du Blind Test
// - Afficher l'interface de jeu (templates/games/blindtest.html)
// - Charger les informations de la salle
// - Vérifier que le joueur appartient à cette salle

// TODO @Quoc Huy: SelectPlaylistHandler - Sélectionner une playlist
// - Méthode POST
// - Récupérer le choix: "Rock", "Rap", ou "Pop"
// - Utiliser le service Spotify pour récupérer les tracks
// - Démarrer le jeu via WebSocket

// TODO @Quoc Huy: SubmitAnswerHandler - Soumettre une réponse
// - Méthode POST
// - Récupérer la réponse du joueur (titre et/ou artiste)
// - Calculer le temps écoulé depuis le début de la track
// - Valider la réponse et attribuer les points (3, 2, 1 ou 0)
// - Envoyer la confirmation via WebSocket
