// RESPONSABLE: @Quoc Huy & @ilian
// Client WebSocket côté frontend

/* TODO @Quoc Huy & @ilian: Implémenter:
   - Connexion au serveur WebSocket
   - Envoi de messages
   - Réception de messages
   - Reconnexion automatique en cas de déconnexion
   - Gestion des erreurs
*/

class GameWebSocket {
    constructor(roomID, userID) {
        this.roomID = roomID;
        this.userID = userID;
        this.ws = null;
        this.handlers = {};
    }
    
    connect() {
        // TODO: Établir la connexion WebSocket
        // this.ws = new WebSocket(`ws://localhost:8080/ws?room=${this.roomID}&user=${this.userID}`);
        
        // this.ws.onopen = () => {
        //     console.log('Connecté au serveur');
        // };
        
        // this.ws.onmessage = (event) => {
        //     const message = JSON.parse(event.data);
        //     this.handleMessage(message);
        // };
    }
    
    send(type, payload) {
        // TODO: Envoyer un message au serveur
        // const message = {
        //     type: type,
        //     room_id: this.roomID,
        //     user_id: this.userID,
        //     payload: payload
        // };
        // this.ws.send(JSON.stringify(message));
    }
    
    on(messageType, handler) {
        // TODO: Enregistrer un handler pour un type de message
        this.handlers[messageType] = handler;
    }
    
    handleMessage(message) {
        // TODO: Router le message vers le bon handler
        if (this.handlers[message.type]) {
            this.handlers[message.type](message);
        }
    }
}
