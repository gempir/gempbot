import { create } from "zustand";

export interface ScTokenContent {
  Login: string;
  UserID: string;
}

export interface Store {
  twitchClientId: string;
  apiBaseUrl: string;
  baseUrl: string;
  yjsWsUrl: string;
  scToken: string | null;
  scTokenContent: ScTokenContent | null;
  setScToken: (token: string) => void;
  setManaging: (managing: string | null) => void;
  managing: string | null;
}

export const useStore = create<Store>((set) => ({
  twitchClientId: "",
  apiBaseUrl: "",
  yjsWsUrl: "",
  baseUrl: "",
  scToken: null,
  scTokenContent: null,
  managing: null,
  setScToken: (token: string) => set({ scToken: token }),
  setManaging: (managing: string | null) => set({ managing }),
}));

// Helper function to initialize store with server state
export const initializeStore = (serverState: Partial<Store>) => {
  useStore.setState(serverState);
};
