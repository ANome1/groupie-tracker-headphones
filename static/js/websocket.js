let socket;

function connectWebSocket(roomCode, playerID) {
    const protocol = window.location.protocol === "https:" ? "wss:" : "ws:";
    const wsUrl = `${protocol}//${window.location.host}/ws?room=${roomCode}&player=${playerID}`;

    socket = new WebSocket(wsUrl);

    socket.onopen = () => {
        console.log("Connected to WebSocket");
        if (window.onWSReady) {
            window.onWSReady();
        }
    };

    socket.onmessage = (event) => {
        const msg = JSON.parse(event.data);
        console.log("Received message:", msg);

        // Vérifier si c'est un message Blind Test (format minuscule)
        if (msg.type) {
            switch (msg.type) {
                case "round_start":
                case "round_end":
                case "scoreboard_update":
                case "game_end":
                case "error":
                    if (window.handleBlindTestMessage) {
                        window.handleBlindTestMessage(msg);
                    }
                    return;
            }
        }

        // Messages Petit Bac (format majuscule Type)
        switch (msg.Type) {
            case "GAME_START":
                // Redirection vers la page de jeu
                window.location.href = msg.Data;
                break;
            case "VALIDATION_PHASE":
                // Géré dans la page de jeu
                if (window.handleValidationPhase) {
                    window.handleValidationPhase(msg.Data);
                }
                break;
            case "SCORES_UPDATE":
                if (window.updateScores) {
                    window.updateScores(msg.Data);
                }
                break;
            case "NEW_ROUND":
                if (window.handleNewRound) {
                    window.handleNewRound(msg.Data);
                }
                break;
            case "GAME_OVER":
                if (window.handleGameOver) {
                    window.handleGameOver(msg.Data);
                }
                break;
            case "PLAYER_JOINED":
                // Mise à jour de la liste des joueurs
                if (window.handlePlayerJoined) {
                    window.handlePlayerJoined(msg.Data);
                }
                break;
            case "PLAYER_LEFT":
                // Mise à jour de la liste des joueurs
                if (window.handlePlayerLeft) {
                    window.handlePlayerLeft(msg.Data);
                }
                break;
            default:
                console.log("Unknown message type:", msg.Type);
        }
    };

    socket.onclose = () => {
        console.log("WebSocket connection closed");
    };

    socket.onerror = (error) => {
        console.error("WebSocket error:", error);
    };
}

// Fonction compatible avec le nouveau format pour blind test
function connectWS(roomCode, userPseudo) {
    connectWebSocket(roomCode, userPseudo);
}

function sendMessage(type, payload) {
    if (socket && socket.readyState === WebSocket.OPEN) {
        socket.send(JSON.stringify({
            Type: type,
            Data: payload
        }));
    } else {
        console.error("WebSocket is not open");
    }
}

// Fonction pour envoyer des messages depuis blindtest.js
function sendWebSocketMessage(message) {
    if (socket && socket.readyState === WebSocket.OPEN) {
        socket.send(JSON.stringify(message));
    } else {
        console.error("WebSocket is not open");
    }
}
