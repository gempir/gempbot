import jwt_decode from "jwt-decode";
import { NextPageContext } from "next";
import { ScTokenContent } from "../store";
import { parseCookie } from "./cookie";

export const initializeStore = (ctx: NextPageContext) => {
    const cookies = parseCookie(ctx.req?.headers.cookie ?? "");

    let scTokenContent;
    if (cookies.scToken) {
        try {
            scTokenContent = jwt_decode<ScTokenContent | undefined>(cookies.scToken ?? "");
        } catch (e) {
            console.error(e);
        }
    }

    return {
        props: {
            store: {
                scTokenContent: scTokenContent ?? "",
                scToken: cookies.scToken ?? "",
                managing: cookies.managing ?? "",
                twitchClientId: (process.env.NEXT_PUBLIC_TWITCH_CLIENT_ID ?? "").replaceAll('"', ''),
                apiBaseUrl: (process.env.NEXT_PUBLIC_API_BASE_URL ?? "").replaceAll('"', ''),
                baseUrl: (process.env.NEXT_PUBLIC_BASE_URL ?? "").replaceAll('"', ''),
            }
        },
    };
}