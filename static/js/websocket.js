let socket;

function connectWebSocket(roomCode, playerID) {
    const protocol = window.location.protocol === "https:" ? "wss:" : "ws:";
    const wsUrl = `${protocol}//${window.location.host}/ws?room=${roomCode}&player=${playerID}`;

    socket = new WebSocket(wsUrl);

    socket.onopen = () => {
        console.log("Connected to WebSocket");
    };

    socket.onmessage = (event) => {
        const msg = JSON.parse(event.data);
        console.log("Received message:", msg);

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

function sendMessage(type, payload) {
    if (socket && socket.readyState === WebSocket.OPEN) {
        // Le serveur attend { Type: string, Data: json.RawMessage }
        // json.RawMessage attend un []byte qui est un JSON valide.
        // Si on envoie un objet JS dans Data, JSON.stringify le convertira en JSON string.
        // Go unmarshalera ce JSON string en []byte.
        socket.send(JSON.stringify({
            Type: type,
            Data: payload // Sera sérialisé comme un objet JSON imbriqué
        }));
    } else {
        console.error("WebSocket is not open");
    }
}
