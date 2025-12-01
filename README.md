# Groupie Tracker - Jeux Multijoueurs Musicaux

Plateforme de jeux musicaux multijoueurs avec **Blind Test** et **Petit Bac**.

## 🎯 Équipe

- **@Nome**: Landing page, Base de données (SQLite), Authentification, Frontend général
- **@Quoc Huy**: Blind Test (Backend + Frontend + Spotify API)
- **@ilian**: Petit Bac (Backend + Frontend + CRUD catégories)

## 🚀 Démarrage Rapide

```bash
# Démarrer le serveur
go run main.go

# Ouvrir dans le navigateur
http://localhost:8080
```

## 📋 Structure du Projet

```
groupie-tracker-headphones/
├── main.go                  # Serveur HTTP (FONCTIONNEL)
├── config/                  # Configuration
├── database/
│   └── schema.sql          # Schéma SQLite (COMPLET)
├── models/                  # Structures de données
├── handlers/                # Handlers HTTP (à implémenter)
├── services/                # Logique métier (à implémenter)
├── websocket/               # WebSocket custom (à implémenter)
├── utils/                   # Utilitaires (à implémenter)
├── templates/               # Pages HTML (FONCTIONNELLES)
└── static/                  # CSS, JS, images
```

## 🎮 Spécifications des Jeux

### Blind Test (@Quoc Huy)
- **Timer**: 37 secondes par musique
- **Playlists**: Rock, Rap, Pop
- **Système de points**:
  - < 10s: 3 points
  - < 20s: 2 points
  - < 37s: 1 point
- **API**: Spotify (OAuth)

### Petit Bac (@ilian)
- **Manches**: 9 (constante `NbrsManche`)
- **Validation**: Collective (2/3 des joueurs)
- **Système de points**:
  - Réponse unique: 2 points
  - Réponse commune: 1 point
- **CRUD**: Catégories personnalisées

## 🎨 Thème Visuel

- **Couleur principale**: Orange `#FF6B35`
- **Couleur secondaire**: Jaune `#FFD23F`

## 🔧 Technologies

- **Backend**: Go 1.21 (100% standard library, AUCUNE dépendance externe)
- **Base de données**: SQLite
- **WebSocket**: Custom (pas de gorilla/websocket)
- **Frontend**: HTML, CSS, JavaScript vanilla


## 📝 État du Projet

✅ **Fonctionnel**:
- Serveur Go avec routes de base
- Templates HTML avec structure complète
- Schéma de base de données SQLite

⏳ **À implémenter** (commentaires TODO dans chaque fichier):
- Connexion à la base de données (@Nome)
- Système d'authentification (@Nome)
- Gestion des salles (@Nome)
- Intégration Spotify API (@Quoc Huy)
- Logique Blind Test (@Quoc Huy)
- Logique Petit Bac (@ilian)
- Hub WebSocket (@Quoc Huy & @ilian)
- Frontend interactif (JavaScript)

