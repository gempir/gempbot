import { NextPageContext } from "next";
import React, { useEffect } from "react";
import { Emotehistory } from "../components/Home/Emotehistory";
import { PredictionLog } from "../components/Home/PredictionLog";
import { Teaser } from "../components/Teaser";
import { parseCookie } from "../service/cookie";
import { Store, useStore } from "../store";

export default function Home({ store }: { store: Store }) {
    useEffect(() => {
        useStore.setState(store);
    }, [store]);

    const isLoggedIn = useStore(s => !!s.scToken);


    return <div className="p-4 w-full max-h-screen flex gap-4">
        {!isLoggedIn && <Teaser />}
        {isLoggedIn && <>
            <Emotehistory />
            <PredictionLog />
        </>}

    </div>
}

Home.getInitialProps = (ctx: NextPageContext) => {
    const cookies = parseCookie(ctx.req?.headers.cookie ?? "");
    useStore.setState({ scToken: cookies.scToken });

    return {
        store: {
            scToken: cookies.scToken ?? "",
        },
    };
}