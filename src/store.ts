import create from 'zustand'

export interface Store {
    twitchClientId: string;
    apiBaseUrl: string;
    baseUrl: string;
    scToken?: string;
    setScToken: (token: string) => void;
    managing?: string;
}

export const useStore = create<Store>(set => ({
    twitchClientId: process.env.REACT_APP_TWITCH_CLIENT_ID ?? "",
    apiBaseUrl: makeUrl(process.env.REACT_APP_BASE_URL ?? process.env.REACT_APP_VERCEL_URL ?? ""),
    baseUrl: makeUrl(process.env.REACT_APP_BASE_URL ?? process.env.REACT_APP_VERCEL_URL ?? ""),
    scToken: undefined,
    setScToken: (token: string) => set(state => ({ scToken: token })),
    managing: undefined,
}));

function isDev() {
    return process.env.REACT_APP_VERCEL_ENV === "development";
}

function makeUrl(domain: string): string {
    if (domain === "") {
        return "";
    }

    const protocol = isDev() ? "http" : "https";

    return `${protocol}://${domain}`;
}