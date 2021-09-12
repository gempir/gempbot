import { Store as PStore } from "pullstate";
import { getCookie } from "./service/cookie";

export interface Store {
    twitchClientId: string,
    baseUrl: string,
    apiBaseUrl: string,
    scToken: string | null,
    managing: string,
}

const url = new URL(window.location.href);

export const store = new PStore<Store>({
    twitchClientId: process.env.REACT_APP_TWITCH_CLIENT_ID ?? "",
    apiBaseUrl: url.protocol + url.host ?? "",
    baseUrl: url.protocol + url.host ?? "",
    scToken: getCookie("scToken"),
    managing: window.localStorage.getItem("managing") ?? "",
});