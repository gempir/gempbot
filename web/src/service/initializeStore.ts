import jwt_decode from "jwt-decode";
import { NextPageContext } from "next";
import { ScTokenContent } from "../store";
import { parseCookie } from "./cookie";

export const initializeStore = (ctx: NextPageContext) => {
    const cookies = parseCookie(ctx.req?.headers.cookie ?? "");

    let scTokenContent = null;
    if (cookies.scToken) {
        try {
            scTokenContent = jwt_decode<ScTokenContent | null>(cookies.scToken ?? "") ?? null;
        } catch (e) {
            console.error(e);
        }
    }

    return {
        props: {
            store: {
                scTokenContent: scTokenContent,
                scToken: cookies.scToken ? cookies.scToken : null,
                managing: cookies.managing ? cookies.managing : null,
                twitchClientId: (process.env.NEXT_PUBLIC_TWITCH_CLIENT_ID ?? "").replaceAll('"', ''),
                apiBaseUrl: (process.env.NEXT_PUBLIC_API_BASE_URL ?? "").replaceAll('"', ''),
                yjsWsUrl: (process.env.NEXT_PUBLIC_YJS_WS_URL ?? "").replaceAll('"', ''),
                baseUrl: (process.env.NEXT_PUBLIC_BASE_URL ?? "").replaceAll('"', ''),
            }
        },
    };
}

export const initializeStoreWithProps = (props: any) => {
    return (ctx: NextPageContext) => {
        const cookies = parseCookie(ctx.req?.headers.cookie ?? "");

        let scTokenContent = null;
        if (cookies.scToken) {
            try {
                scTokenContent = jwt_decode<ScTokenContent | null>(cookies.scToken ?? "") ?? null;
            } catch (e) {
                console.error(e);
            }
        }

        return {
            props: {
                ...props,
                store: {
                    scTokenContent: scTokenContent,
                    scToken: cookies.scToken ? cookies.scToken : null,
                    managing: cookies.managing ? cookies.managing : null,
                    twitchClientId: (process.env.NEXT_PUBLIC_TWITCH_CLIENT_ID ?? "").replaceAll('"', ''),
                    apiBaseUrl: (process.env.NEXT_PUBLIC_API_BASE_URL ?? "").replaceAll('"', ''),
                    yjsWsUrl: (process.env.NEXT_PUBLIC_YJS_WS_URL ?? "").replaceAll('"', ''),
                    baseUrl: (process.env.NEXT_PUBLIC_BASE_URL ?? "").replaceAll('"', ''),
                }
            },
        };
    }
}