-- RESPONSABLE: @Nome
-- Schéma de base de données SQLite pour Groupie Tracker
-- Ce fichier définit toutes les tables nécessaires au projet

-- Table des utilisateurs
-- Stocke les informations d'authentification
CREATE TABLE IF NOT EXISTS users (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    username TEXT UNIQUE NOT NULL,
    email TEXT UNIQUE NOT NULL,
    password_hash TEXT NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- Table des salles de jeu
-- Une salle peut héberger un Blind Test ou un Petit Bac
CREATE TABLE IF NOT EXISTS rooms (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL,
    code TEXT UNIQUE NOT NULL, -- Code à 6 caractères pour rejoindre
    host_id INTEGER NOT NULL,
    game_type TEXT NOT NULL, -- 'blindtest' ou 'petitbac'
    max_players INTEGER DEFAULT 8,
    status TEXT DEFAULT 'waiting', -- 'waiting', 'in_progress', 'finished'
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (host_id) REFERENCES users(id)
);

-- Table des participants dans les salles
CREATE TABLE IF NOT EXISTS room_participants (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    room_id INTEGER NOT NULL,
    user_id INTEGER NOT NULL,
    score INTEGER DEFAULT 0,
    joined_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (room_id) REFERENCES rooms(id),
    FOREIGN KEY (user_id) REFERENCES users(id)
);

-- Table des sessions de jeu (historique)
CREATE TABLE IF NOT EXISTS game_sessions (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    room_id INTEGER NOT NULL,
    game_data TEXT, -- JSON pour stocker l'état du jeu
    started_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    ended_at DATETIME,
    FOREIGN KEY (room_id) REFERENCES rooms(id)
);

-- TODO @ilian: Table pour les catégories personnalisées du Petit Bac
-- CREATE TABLE IF NOT EXISTS petitbac_categories (
--     id INTEGER PRIMARY KEY AUTOINCREMENT,
--     name TEXT NOT NULL,
--     created_by INTEGER,
--     is_default BOOLEAN DEFAULT 0,
--     FOREIGN KEY (created_by) REFERENCES users(id)
-- );

-- Index pour améliorer les performances
CREATE INDEX IF NOT EXISTS idx_rooms_code ON rooms(code);
CREATE INDEX IF NOT EXISTS idx_room_participants ON room_participants(room_id, user_id);
