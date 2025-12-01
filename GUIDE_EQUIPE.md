# 📘 Guide pour l'Équipe Groupie Tracker

## 🎯 Vue d'ensemble

Tous les fichiers Go, JavaScript et CSS sont **vides** avec uniquement des **commentaires TODO** qui expliquent ce que chaque membre doit implémenter.

## 👥 Répartition des Tâches

### @Nome - Infrastructure & Authentification
**Fichiers à implémenter:**
- `config/config.go` - Configuration de l'application
- `models/user.go` - Modèle utilisateur
- `models/room.go` - Modèle de salle
- `handlers/auth.go` - Inscription, connexion, déconnexion
- `handlers/room.go` - Création, rejoindre, lobby
- `services/auth_service.go` - Logique d'authentification
- `services/room_service.go` - Logique des salles
- `utils/validator.go` - Validation des données
- `utils/hash.go` - Hashage SHA256
- `utils/helpers.go` - Fonctions utilitaires
- `static/css/style.css` - Styles généraux
- `static/css/auth.css` - Styles authentification
- `static/css/responsive.css` - Responsive design
- `static/js/animations.js` - Animations UI
- `templates/landing.html` - Page d'accueil (déjà créée)

**Priorités:**
1. Connexion SQLite dans `main.go` (ligne 14)
2. Système d'authentification complet
3. Gestion des salles
4. CSS et design général

### @Quoc Huy - Blind Test & Spotify
**Fichiers à implémenter:**
- `models/blindtest.go` - Modèle du jeu Blind Test
- `handlers/blindtest.go` - Handlers du Blind Test
- `services/spotify_service.go` - Intégration Spotify API
- `websocket/hub.go` - Hub WebSocket (avec @ilian)
- `websocket/client.go` - Client WebSocket (avec @ilian)
- `websocket/messages.go` - Messages WebSocket (avec @ilian)
- `static/css/game.css` - Styles Blind Test (partie)
- `static/js/websocket.js` - WebSocket client (avec @ilian)
- `static/js/blindtest.js` - Interface Blind Test
- `templates/games/blindtest.html` - Page Blind Test (déjà créée)

**Priorités:**
1. Compte développeur Spotify et OAuth
2. Récupération des tracks depuis l'API
3. WebSocket pour la communication en temps réel
4. Système de timer (37s) et de points (3/2/1)
5. Interface JavaScript pour le jeu

**Important:**
- Timer: **37 secondes** par musique (constante `DefaultBlindTestTimer`)
- Playlists: Rock, Rap, Pop
- Système de points basé sur la rapidité

### @ilian - Petit Bac & Catégories
**Fichiers à implémenter:**
- `models/petitbac.go` - Modèle du jeu Petit Bac
- `handlers/petitbac.go` - Handlers du Petit Bac + CRUD catégories
- `websocket/hub.go` - Hub WebSocket (avec @Quoc Huy)
- `websocket/client.go` - Client WebSocket (avec @Quoc Huy)
- `websocket/messages.go` - Messages WebSocket (avec @Quoc Huy)
- `static/css/game.css` - Styles Petit Bac (partie)
- `static/js/websocket.js` - WebSocket client (avec @Quoc Huy)
- `static/js/petitbac.js` - Interface Petit Bac
- `templates/games/petitbac.html` - Page Petit Bac (déjà créée)

**Priorités:**
1. WebSocket pour la communication en temps réel
2. Système de validation collective (2/3 des joueurs)
3. Calcul des points (unique vs commun)
4. CRUD des catégories personnalisées
5. Interface JavaScript pour le jeu

**Important:**
- Nombre de manches: **9** (constante `NbrsManche`)
- Variable `scoreboardActualPointInGame` pour les points de la manche actuelle
- Validation: 2/3 des joueurs doivent approuver

## 🚀 Comment Démarrer

### 1. Cloner et Tester le Serveur
```bash
cd /home/anome/B1/groupie-tracker-headphones
go run main.go
```
Ouvrir http://localhost:8080 pour voir l'aperçu

### 2. Lire les Commentaires TODO
Chaque fichier contient des commentaires détaillés expliquant:
- Les structures à créer
- Les fonctions à implémenter
- Les étapes d'implémentation
- Les spécifications techniques

### 3. Implémenter par Étapes
Ne pas essayer de tout faire d'un coup ! Suivre cet ordre:

**Phase 1 - Base (@Nome)**
1. Connexion SQLite dans `main.go`
2. Modèles `User` et `Room`
3. Handlers d'authentification
4. CSS de base

**Phase 2 - Jeux (@Quoc Huy & @ilian)**
1. Modèles des jeux
2. Hub WebSocket
3. Handlers des jeux
4. Services (Spotify, Game)

**Phase 3 - Frontend**
1. JavaScript WebSocket
2. Interface des jeux
3. Animations et CSS avancé

## 📝 Conventions de Code

### Go
```go
// TODO @NomDuMembre: Description de la tâche
// - Étape 1
// - Étape 2
```

### Commits
```bash
git add .
git commit -m "@NomDuMembre: Description courte de ce qui a été fait"
git push
```

Exemples:
- `@Nome: Ajout de la connexion SQLite et authentification`
- `@Quoc Huy: Intégration Spotify API et handlers Blind Test`
- `@ilian: Système de validation collective pour Petit Bac`

## 🔧 Contraintes Techniques

### ❌ INTERDIT
- Utiliser des dépendances GitHub (gorilla/websocket, bcrypt, etc.)
- JavaScript côté backend (uniquement pour le navigateur)
- Copier-coller du code sans comprendre

### ✅ AUTORISÉ
- Packages Go standard uniquement (`net/http`, `database/sql`, `crypto/sha256`, etc.)
- JavaScript vanilla dans le navigateur
- SQLite (avec driver `github.com/mattn/go-sqlite3` si nécessaire pour le driver SQL)

## 🎮 Spécifications Détaillées

### Blind Test
```
Durée: 37 secondes par musique
Playlists: Rock, Rap, Pop (Spotify)
Points:
  - < 10s: 3 points
  - < 20s: 2 points
  - < 37s: 1 point
  - Après 37s ou faux: 0 point
```

### Petit Bac
```
Manches: 9 (NbrsManche)
Catégories: Défaut + Personnalisées (CRUD)
Validation: 2/3 des joueurs minimum
Points:
  - Réponse unique: 2 points
  - Réponse commune: 1 point
  - Réponse invalidée: 0 point
Variable: scoreboardActualPointInGame (points de la manche)
```

## 🆘 En Cas de Problème

1. **Relire les commentaires TODO** dans le fichier concerné
2. **Consulter le schéma SQL** dans `database/schema.sql`
3. **Vérifier les modèles** dans `models/` pour comprendre les structures
4. **Tester régulièrement** avec `go run main.go`
5. **Communiquer avec l'équipe** si bloqué

## 📚 Ressources Utiles

- **Spotify API**: https://developer.spotify.com/documentation/web-api
- **WebSocket Go**: https://pkg.go.dev/golang.org/x/net/websocket (ou custom)
- **SQLite Go**: https://github.com/mattn/go-sqlite3
- **SHA256**: https://pkg.go.dev/crypto/sha256

## ✅ Checklist par Membre

### @Nome
- [ ] Connexion SQLite fonctionnelle
- [ ] Système d'authentification complet
- [ ] Gestion des salles (création, rejoindre, quitter)
- [ ] CSS avec thème orange (#FF6B35) et jaune (#FFD23F)
- [ ] Responsive design
- [ ] Landing page attractive

### @Quoc Huy
- [ ] Compte Spotify Developer créé
- [ ] OAuth Spotify fonctionnel
- [ ] Récupération de tracks depuis les playlists
- [ ] Hub WebSocket opérationnel
- [ ] Timer de 37 secondes
- [ ] Système de points (3/2/1)
- [ ] Interface Blind Test complète

### @ilian
- [ ] Hub WebSocket opérationnel
- [ ] Système de validation 2/3 joueurs
- [ ] Calcul points (unique vs commun)
- [ ] CRUD catégories personnalisées
- [ ] 9 manches par partie
- [ ] scoreboardActualPointInGame mis à jour
- [ ] Interface Petit Bac complète

---


