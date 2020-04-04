export default class EventService {
    constructor(apiBaseUrl, callback) {
        this.onEvent = callback;

        const socket = new WebSocket(`${apiBaseUrl.replace("https://", "wss://").replace("http://", "ws://")}/api/ws`);

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