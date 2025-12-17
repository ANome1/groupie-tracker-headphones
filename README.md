# 🎵 Groupie Tracker - Jeux Multijoueurs Musicaux

Plateforme de jeux musicaux multijoueurs complète avec **Blind Test** et **Petit Bac** en temps réel via WebSocket.

## 🎯 Équipe

- **@Nome**: Landing page, Base de données (SQLite), Authentification, Gestion des salles, frontend général
- **@Quoc Huy**: Blind Test (Backend + API Deezer ) websocket
- **@ilian**: Petit Bac (Backend) websocket

## 🚀 Démarrage Rapide

```bash
# Démarrer le serveur
go run server.go

# Ouvrir dans le navigateur
http://localhost:8080
```

## 📋 Architecture du Projet

```
groupie-tracker-headphones/
├── server.go                    # Serveur HTTP (✅ Production)
├── go.mod                       # Dépendances Go
│
├── config/
│   └── config.go               # Configuration de l'application
│
├── database/
│   ├── db.go                   # Connexion et gestion SQLite
│   └── schema.sql              # Schéma complet de la BD
│
├── models/
│   ├── user.go                 # Utilisateurs & authentification
│   ├── room.go                 # Gestion des salles
│   ├── blindtest.go            # État du jeu Blind Test
│   ├── petitbac.go             # État du jeu Petit Bac
│   └── game.go                 # Structures communes
│
├── handlers/
│   ├── auth.go                 # Routes auth (login, register)
│   ├── blindtest.go            # Page Blind Test
│   ├── room.go                 # Gestion des salles (create, join, lobby)
│   ├── petitbac.go             # Page Petit Bac
│   └── websocket.go            # Upgrade WebSocket
│
├── services/
│   ├── auth_service.go         # Logique authentification
│   ├── room_service.go         # Logique gestion des salles
│   ├── blindtest_service.go    # Logique métier Blind Test
│   ├── game_service.go         # Utilitaires jeux
│   ├── deezer_service.go       # API Deezer (musiques)
│   └── spotify_service.go      # API Spotify (alternative)
│
├── websocket/
│   ├── hub.go                  # Hub central (broadcast par salle)
│   ├── client.go               # Client WebSocket (read/write pumps)
│   └── messages.go             # Types de messages
│
├── utils/
│   ├── hash.go                 # Hachage des mots de passe
│   ├── helpers.go              # Utilitaires généraux
│   ├── session.go              # Gestion des cookies
│   └── validator.go            # Validation des entrées
│
├── templates/
│   ├── home.html               # Page d'accueil
│   ├── landing.html            # Landing page
│   ├── auth/
│   │   ├── login.html          # Connexion
│   │   └── register.html       # Inscription
│   ├── games/
│   │   ├── blindtest.html      # Interface Blind Test
│   │   └── petitbac.html       # Interface Petit Bac
│   ├── room/
│   │   ├── create.html         # Créer une salle
│   │   ├── join.html           # Rejoindre une salle
│   │   └── lobby.html          # Lobby d'attente
│   └── components/
│       ├── header.html         # En-tête commun
│       └── footer.html         # Pied de page commun
│
└── static/
    ├── css/
    │   ├── main.css            # Styles globaux
    │   ├── auth.css            # Styles authentification
    │   ├── game.css            # Styles des jeux
    │   └── responsive.css      # Design responsive
    ├── js/
    │   ├── websocket.js        # Gestion WebSocket client
    │   ├── blindtest.js        # Logique client Blind Test
    │   ├── petitbac.js         # Logique client Petit Bac
    │   └── animations.js       # Animations CSS
    └── assets/
        ├── fonts/              # Polices personnalisées
        ├── images/
        │   └── icons/          # Icônes
        └── sounds/             # Effets sonores
```

## 🎮 Spécifications des Jeux

### 🎧 Blind Test
- **Durée par manche**: (configurable)
- **Nombre de manches**: (configurable)
- **Playlists Deezer**: Rock, Pop, Rap US, Rap FR, Indie
- **Mécanisme de scoring**:
  - Réponse complète (titre + artiste) : **2x points**
  - Réponse partielle (titre OU artiste) : **1x point**
- **Synchronisation WebSocket**: Temps restant actualisé en temps réel
- **Révélation progressive**: Blur de la cover enlevé à 1/3 du temps

### 📝 Petit Bac
- **Nombre de manches**: (configurable)
- **Catégories**: Artiste, Groupe de musique, Album, Instrument, Featuring
- **Durée par manche**: (configurable)
- **Phase de validation**: Vote des autres joueurs (2/3+ pour valider)
- **Système de points**:
  - Réponse unique (validée) : **2 points**
  - Réponse commune (validée) : **1 point**
  - Réponse invalide : **0 point**

## 🔧 Technologies

- **Backend**: Go 1.21 (100% standard library)
- **Base de données**: SQLite3
- **WebSocket**: Custom (sans dépendance externe)
- **Frontend**: HTML5, CSS3, JavaScript vanilla
- **API Musicale**: Deezer (proxy local pour CORS)
- **Authentification**: Cookies + Hachage bcrypt

## ✅ État du Projet

### 🏁 Implémenté et Fonctionnel

**Backend (Go)**:
- ✅ Serveur HTTP avec routage complet
- ✅ Authentification (inscription, connexion, cookies sécurisés)
- ✅ Gestion des salles (création, adhésion, lobby)
- ✅ WebSocket Hub avec broadcast par salle
- ✅ Logique Blind Test (timer, validation fuzzy, scoring 2x)
- ✅ Logique Petit Bac (phases, validation, scoring)
- ✅ Base de données SQLite (migrations, schéma)
- ✅ Intégration API Deezer (playlists, preview audio)
- ✅ Synchronisation en temps réel entre clients

**Frontend (JavaScript)**:
- ✅ Interface Blind Test (timer, révélation blur, scores)
- ✅ Interface Petit Bac (vote, validation, résultats)
- ✅ Gestion des salles (create, join, lobby)
- ✅ Authentification (login, register)
- ✅ Synchronisation WebSocket (messages, routage)
- ✅ Animations fluides et responsive design
- ✅ Gestion des erreurs et déconnexions


## 📊 Caractéristiques

- **Multijoueurs temps réel** via WebSocket
- **Gestion des salles** avec code d'accès unique
- **Authentification sécurisée** avec hachage des mots de passe
- **Deux jeux distincts** avec règles et interfaces différentes
- **Scoring dynamique** synchronisé entre tous les clients
- **Zéro dépendance externe** (Go standard library)


