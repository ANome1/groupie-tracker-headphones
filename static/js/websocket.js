let socket;

// Initialisation et routage des messages WebSocket depuis le serveur
function connectWebSocket(roomCode, playerID) {
    const protocol = window.location.protocol === "https:" ? "wss:" : "ws:";
    const wsUrl = `${protocol}//${window.location.host}/ws?room=${roomCode}&player=${playerID}`;

    socket = new WebSocket(wsUrl);

    socket.onopen = () => {
        if (window.onWSReady) {
            window.onWSReady();
        }
    };

    socket.onmessage = (event) => {
        const lines = event.data.split('\n').filter(line => line.trim().length > 0);
        
        for (const line of lines) {
            try {
                const msg = JSON.parse(line);

                // Synchronisation Blind Test: mesages minuscules (round_start, round_end, etc.)
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
                            continue;
                    }
                }

                // Routage Petit Bac: messages majuscules (VALIDATION_PHASE, SCORES_UPDATE, etc.)
                switch (msg.Type) {
                    case "GAME_START":
                        window.location.href = msg.Data;
                        break;
                    case "VALIDATION_PHASE":
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
                    case "ROUND_RESULTS":
                        if (window.handleRoundResults) {
                            window.handleRoundResults(msg.Data);
                        }
                        break;
                    case "GAME_OVER":
                        if (window.handleGameOver) {
                            window.handleGameOver(msg.Data);
                        }
                        break;
                    case "PLAYER_JOINED":
                        if (window.handlePlayerJoined) {
                            window.handlePlayerJoined(msg.Data);
                        }
                        break;
                    case "PLAYER_LEFT":
                        if (window.handlePlayerLeft) {
                            window.handlePlayerLeft(msg.Data);
                        }
                        break;
                }
            } catch (e) {
            }
        }
    };

    socket.onclose = () => {
    };

    socket.onerror = (error) => {
    };
}

function connectWS(roomCode, userPseudo) {
    connectWebSocket(roomCode, userPseudo);
}

function sendMessage(type, payload) {
    if (socket && socket.readyState === WebSocket.OPEN) {
        socket.send(JSON.stringify({
            Type: type,
            Data: payload
        }));
    }
}

function sendWebSocketMessage(message) {
    if (socket && socket.readyState === WebSocket.OPEN) {
        socket.send(JSON.stringify(message));
    }
}
