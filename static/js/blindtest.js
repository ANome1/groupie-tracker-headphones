// RESPONSABLE: @Quoc Huy
// Logique côté client du Blind Test

/* TODO @Quoc Huy: Implémenter:
   
   1. Sélection de playlist:
      - Gérer le clic sur Rock, Rap, ou Pop
      - Envoyer le choix au serveur
      - Cacher la sélection et afficher le jeu
   
   2. Gestion de l'audio:
      - Lecture automatique de l'extrait
      - Arrêt après 37 secondes
   
   3. Timer:
      - Décompte de 37 à 0
      - Animation (changement de couleur < 10s)
      - Calcul du temps écoulé pour les points
   
   4. Soumission de réponse:
      - Capturer la réponse du joueur
      - Calculer le temps écoulé
      - Envoyer via WebSocket
   
   5. Réception des événements:
      - Nouvelle musique (blindtest_track)
      - Mise à jour des scores (score_update)
      - Fin du jeu (game_end)
   
   6. Affichage:
      - Mettre à jour le scoreboard en temps réel
      - Afficher le numéro du round
      - Animations de points gagnés
*/

let currentGame = null;
let startTime = null;

function selectPlaylist(playlist) {
    // TODO: Envoyer le choix de playlist au serveur
    console.log('Playlist sélectionnée:', playlist);
}

function submitAnswer() {
    // TODO: Soumettre la réponse
    const answer = document.getElementById('answer-input').value;
    const timeElapsed = Math.floor((Date.now() - startTime) / 1000);
    console.log('Réponse:', answer, 'Temps:', timeElapsed);
}

function updateTimer(seconds) {
    // TODO: Mettre à jour l'affichage du timer
    document.getElementById('timer').textContent = seconds;
}

function updateScoreboard(scores) {
    // TODO: Mettre à jour le scoreboard
    console.log('Scores:', scores);
}
