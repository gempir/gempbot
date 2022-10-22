import create from 'zustand'

export interface Store {
    twitchClientId: string;
    apiBaseUrl: string;
    baseUrl: string;
    scToken?: string;
    setScToken: (token: string) => void;
    setManaging: (managing: string) => void;
    managing?: string;
}


export const useStore = create<Store>(set => ({
    twitchClientId: (process.env.NEXT_PUBLIC_TWITCH_CLIENT_ID ?? "").replaceAll('"', ''),
    apiBaseUrl: (process.env.NEXT_PUBLIC_API_BASE_URL ?? "").replaceAll('"', ''),
    baseUrl: (process.env.NEXT_PUBLIC_BASE_URL ?? "").replaceAll('"', ''),
    scToken: undefined,
    setScToken: (token: string) => set(state => ({ scToken: token })),
    setManaging: (managing: string) => set(state => ({ managing: managing })),
    managing: undefined,
}));