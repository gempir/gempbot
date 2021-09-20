import { NextPageContext } from "next";
import { parseCookie } from "./cookie";

export const initializeStore = (ctx: NextPageContext) => {
    const cookies = parseCookie(ctx.req?.headers.cookie ?? "");

    return {
        store: {
            scToken: cookies.scToken ?? undefined,
            managing: cookies.managing ?? undefined,
        },
    };
}