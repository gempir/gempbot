import useWebSocket from "react-use-websocket";
import { useStore } from "../store";

export enum WsAction {
    PLAYER_STATE = "PLAYER_STATE",
    JOIN = "JOIN",
    DEBUG = "DEBUG",
}

export function useWs(onMessage = (event: MessageEvent<any>) => { }) {
    const wsApiBaseUrl = useStore(store => store.apiBaseUrl).replace('https://', 'wss://').replace('http://', 'ws://');

    return useWebSocket(wsApiBaseUrl + "/api/ws", {
        onMessage: onMessage
    });
}