package main

// RESPONSABLE: @Nome (infrastructure), @Quoc Huy (WebSocket Blind Test), @ilian (WebSocket Petit Bac)
// Ce fichier initialise le serveur web, les routes et les connexions

func main() {
	// TODO @Nome: Initialiser la connexion à la base de données SQLite
	// - Charger le fichier database/groupie-tracker.db
	// - Exécuter le schéma SQL si la base n'existe pas

	// TODO @Nome: Configurer les routes HTTP
	// - Route / (landing page)
	// - Routes /register et /login (authentification)
	// - Routes /room/* (création/rejoindre salles)

	// TODO @Quoc Huy: Ajouter les routes du Blind Test
	// - Route /game/blindtest (interface de jeu)
	// - Route /api/blindtest/* (endpoints API)

	// TODO @ilian: Ajouter les routes du Petit Bac
	// - Route /game/petitbac (interface de jeu)
	// - Route /api/petitbac/* (endpoints API)
	// - Routes pour gérer les catégories personnalisées (CRUD)

	// TODO @Quoc Huy & @ilian: Initialiser le hub WebSocket
	// - Lancer le hub.Run() en goroutine
	// - Route /ws pour les connexions WebSocket

	// TODO @Nome: Servir les fichiers statiques
	// - CSS, JS, images depuis le dossier static/

	// TODO: Démarrer le serveur HTTP sur le port 8080
	// http.ListenAndServe(":8080", nil)
}
