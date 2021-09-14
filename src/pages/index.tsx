import { NextPageContext } from "next";
import React from "react";
import { Teaser } from "../components/Teaser";
import { parseCookie } from "../service/cookie";
import { useStore } from "../store";

export default function Home() {
    console.log("xd");
    return <div className="p-4 w-full max-h-screen flex gap-4">
        <Teaser />
        {/* <Emotehistory />
        <PredictionLog /> */}
    </div>
}

export async function getServerSideProps(ctx: NextPageContext) {
    const cookies = parseCookie(ctx.req?.headers.cookie ?? "");
    useStore.setState({ scToken: cookies.scToken });

    return {
        props: {
            scToken: cookies.scToken ?? "",
        },
    };
}