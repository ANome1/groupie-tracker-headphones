let socket;

document.addEventListener('DOMContentLoaded', () => {
    initWebSocket();
});

function initWebSocket() {
    // Determine protocol (ws or wss) based on current page
    const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:';
    const wsUrl = `${protocol}//${window.location.host}/ws`;
    
    console.log("Connecting to WebSocket at", wsUrl);
    socket = new WebSocket(wsUrl);

    socket.onopen = () => {
        console.log("Connected to WebSocket");
    };

    socket.onmessage = (event) => {
        try {
            const msg = JSON.parse(event.data);
            console.log("Received:", msg);
            // TODO: Handle server messages (update UI, etc.)
        } catch (e) {
            console.error("Error parsing message:", e);
        }
    };

    socket.onclose = () => {
        console.log("Disconnected from WebSocket");
    };
    
    socket.onerror = (err) => {
        console.error("WebSocket error:", err);
    };
}

// Generic function to send messages to the server
function sendSocketMessage(type, data) {
    if (socket && socket.readyState === WebSocket.OPEN) {
        socket.send(JSON.stringify({
            type: type,
            data: data
        }));
    } else {
        console.error("WebSocket is not connected");
    }
}
