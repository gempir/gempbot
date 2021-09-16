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
    twitchClientId: process.env.NEXT_PUBLIC_TWITCH_CLIENT_ID ?? "",
    apiBaseUrl: makeUrl(process.env.NEXT_PUBLIC_BASE_URL ?? process.env.NEXT_PUBLIC_VERCEL_URL ?? ""),
    baseUrl: makeUrl(process.env.NEXT_PUBLIC_BASE_URL ?? process.env.NEXT_PUBLIC_VERCEL_URL ?? ""),
    scToken: undefined,
    setScToken: (token: string) => set(state => ({ scToken: token })),
    setManaging: (managing: string) => set(state => ({ managing: managing })),
    managing: undefined,
}));

function isDev() {
    return process.env.NEXT_PUBLIC_VERCEL_ENV === "development";
}

function makeUrl(domain: string): string {
    if (domain === "") {
        return "";
    }

    const protocol = isDev() ? "http" : "https";

    return `${protocol}://${domain}`;
}