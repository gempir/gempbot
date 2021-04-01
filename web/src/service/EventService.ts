export default class EventService {
    constructor(apiBaseUrl: string, callback: (data: Record<string, unknown>) => void) {

        function connect() {
            var ws = new WebSocket(`${apiBaseUrl.replace("https://", "wss://").replace("http://", "ws://")}/api/ws`);
            
            ws.onmessage = (event) => {
                callback(JSON.parse(event.data));
            };

            ws.onclose = e => {
                console.log('Socket is closed. Reconnect will be attempted in 1 second.', e.reason);
                setTimeout(connect, 1000);
            };

            ws.onerror = err => {
                console.error('Socket encountered error: ', err, 'Closing socket');
                ws.close();
            };
        }

        connect();
    }
}

