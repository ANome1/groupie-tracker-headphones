// RESPONSABLE: @ilian
// Logique côté client du Petit Bac

/* TODO @ilian: Implémenter:
   
   1. Affichage de la lettre:
      - Afficher la lettre tirée en grand
      - Animation d'apparition
   
   2. Gestion des catégories:
      - Générer dynamiquement les champs de saisie
      - Une ligne par catégorie
      - Focus automatique sur le premier champ
   
   3. CRUD des catégories personnalisées:
      - Modal pour ajouter/modifier/supprimer
      - Requêtes AJAX vers l'API
      - Mise à jour de la liste
   
   4. Timer:
      - Décompte (généralement 60 secondes)
      - Soumission automatique à la fin
   
   5. Soumission des réponses:
      - Récupérer toutes les réponses
      - Vérifier qu'elles commencent par la bonne lettre
      - Envoyer via WebSocket
   
   6. Phase de validation:
      - Afficher toutes les réponses de tous les joueurs
      - Permettre de voter (valider/invalider)
      - Système: 2/3 des joueurs doivent valider
      - Envoyer les votes via WebSocket
   
   7. Réception des événements:
      - Nouvelle lettre (petitbac_letter)
      - Résultats de validation (petitbac_validate)
      - Mise à jour des scores (score_update)
      - scoreboardActualPointInGame pour temps réel
   
   8. Affichage:
      - Mettre à jour le scoreboard
      - Afficher le numéro de la manche (sur NbrsManche = 9)
      - Animations de points (2 pts unique, 1 pt commun)
*/

let categories = [];
const NbrsManche = 9; // Nombre total de manches

function generateCategoriesGrid() {
    // TODO: Générer les champs de saisie pour chaque catégorie
    const container = document.getElementById('categories');
    categories.forEach(category => {
        // Créer un champ pour chaque catégorie
    });
}

function submitAnswers() {
    // TODO: Récupérer et envoyer toutes les réponses
    const answers = {};
    // Collecter les réponses...
    console.log('Réponses:', answers);
}

function displayValidationPhase(allPlayersAnswers) {
    // TODO: Afficher toutes les réponses pour validation collective
    console.log('Validation:', allPlayersAnswers);
}

function voteForAnswer(userID, category, isValid) {
    // TODO: Envoyer un vote de validation
    console.log('Vote:', userID, category, isValid);
}

function manageCategoriesModal() {
    // TODO: Ouvrir le modal de gestion des catégories
    document.getElementById('categories-modal').style.display = 'block';
}

function loadCategories() {
    // TODO: Charger les catégories depuis l'API
    // fetch('/api/petitbac/categories')
}

function addCategory(name) {
    // TODO: Ajouter une catégorie personnalisée
    // fetch('/api/petitbac/categories', { method: 'POST', ... })
}
