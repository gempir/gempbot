import useWebSocket from "react-use-websocket";
import { useStore } from "../store";

export enum WsAction {
  PLAYER_STATE = "PLAYER_STATE",
  JOIN = "JOIN",
  GET_QUEUE = "GET_QUEUE",
  QUEUE_STATE = "QUEUE_STATE",
  DEBUG = "DEBUG",
}

export function useWs(onMessage = (_event: MessageEvent<any>) => {}) {
  const wsApiBaseUrl = useStore((store) => store.apiBaseUrl)
    .replace("https://", "wss://")
    .replace("http://", "ws://");

  return useWebSocket(`${wsApiBaseUrl}/api/ws`, {
    onMessage: onMessage,
    shouldReconnect: (_closeEvent) => true,
    reconnectAttempts: 10,
    reconnectInterval: 3000,
  });
}
