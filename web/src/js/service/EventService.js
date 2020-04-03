export default class EventService {
    constructor(callback) {
        this.onEvent = callback;

        const socket = new WebSocket("ws://localhost:8000/api/ws");

        socket.onopen = (e) => {
            console.log("[open] Connection established");
        };

        socket.onmessage = (event) => {
            this.onEvent(JSON.parse(event.data));
        };

        socket.onclose = (event) => {
            if (event.wasClean) {
                console.log(`[close] Connection closed cleanly, code=${event.code} reason=${event.reason}`);
            } else {
                console.log('[close] Connection died');
            }
        };

        socket.onerror = (error) => {
            console.error(error);
        };
    }
}