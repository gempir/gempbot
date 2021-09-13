import { Store as PStore } from "pullstate";
import { getCookie } from "./service/cookie";

export interface Store {
    twitchClientId: string,
    baseUrl: string,
    apiBaseUrl: string,
    scToken: string | null,
    managing: string,
}

function isDev() {
    return process.env.REACT_APP_VERCEL_ENV === "development";
}

const env = process.env;

export const store = new PStore<Store>({
    twitchClientId: process.env.REACT_APP_TWITCH_CLIENT_ID ?? "",
    apiBaseUrl: makeUrl(env.REACT_APP_BASE_URL ?? env.REACT_APP_VERCEL_URL ?? ""),
    baseUrl: makeUrl(env.REACT_APP_BASE_URL ?? env.REACT_APP_VERCEL_URL ?? ""),
    scToken: getCookie("scToken"),
    managing: "",
    // managing: window.localStorage.getItem("managing") ?? "",
});

function makeUrl(domain: string): string {
    if (domain === "") {
        return "";
    }

    const protocol = isDev() ? "http" : "https";

    return `${protocol}://${domain}`;
}