import useWebSocket from "react-use-websocket";
import { useStore } from "../store";

export function useWs() {
    const wsApiBaseUrl = useStore(store => store.apiBaseUrl).replace('https://', 'wss://').replace('http://', 'ws://');

    return useWebSocket(wsApiBaseUrl + "/api/ws");
}