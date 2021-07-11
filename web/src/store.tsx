import { Store as PStore } from "pullstate";
import { getCookie } from "./service/cookie";

export interface Store {
    twitchClientId: string,
    baseUrl: string,
    apiBaseUrl: string,
    scToken: string | null,
    managing: string,
}

export const store = new PStore<Store>({
    twitchClientId: process.env.REACT_APP_TWITCH_CLIENT_ID ?? "",
    apiBaseUrl: process.env.REACT_APP_API_BASE_URL ?? "",
    baseUrl: process.env.REACT_APP_BASE_URL ?? "",
    scToken: getCookie("scToken"),
    managing: window.localStorage.getItem("managing") ?? "",
});