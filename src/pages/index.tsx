import { NextPageContext } from "next";
import React from "react";
// import { parseCookie } from "../service/parseCookie";

export default function Home() {
    return <div className="p-4 w-full max-h-screen flex gap-4">
        {/* <Emotehistory />
        <PredictionLog /> */}
    </div>
}

export async function getServerSideProps(ctx: NextPageContext) {
    // const cookies = parseCookie(ctx.req?.headers.cookie ?? "");
    const cookies = { scToken: "" };

    return {
        props: {
            scToken: cookies.scToken ?? "",
        },
    };
}