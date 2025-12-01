package websocket

// RESPONSABLE: @Quoc Huy & @ilian
// Types de messages WebSocket

// TODO @Quoc Huy & @ilian: Définir les types de messages:
// - Type constants: "player_joined", "player_left", "game_start", etc.
// - Structure Message avec: Type, Payload, RoomID, UserID

// TODO @Quoc Huy: Messages spécifiques au Blind Test
// - "track_start": Nouvelle musique commence
// - "answer_submit": Joueur soumet une réponse
// - "answer_result": Résultat de la réponse (correct/incorrect, points)
// - "round_end": Fin du round avec le classement
// - "game_end": Fin de la partie avec le classement final

// TODO @ilian: Messages spécifiques au Petit Bac
// - "round_start": Nouvelle manche avec la lettre
// - "answers_submit": Joueur soumet ses réponses
// - "validation_phase": Phase de validation collective
// - "validation_vote": Vote d'un joueur sur les réponses
// - "round_results": Résultats de la manche (points attribués)
// - "game_end": Fin de la partie avec le classement final
