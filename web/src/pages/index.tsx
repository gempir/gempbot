import React, { useEffect } from "react";
import { Emotehistory } from "../components/Home/Emotehistory";
import { Teaser } from "../components/Teaser";
import { initializeStore } from "../service/initializeStore";
import { useStore } from "../store";

export default function Home() {
    const isLoggedIn = useStore(s => !!s.scToken);

    useEffect(() => {
        const path = window.localStorage.getItem("redirect");

        if (path) {
            window.localStorage.removeItem("redirect");
            window.location.href = path;
        }
    }, []);

    return <div className="p-4 w-full max-h-screen flex gap-4">
        {!isLoggedIn && <Teaser />}
        {isLoggedIn && <>
            <Emotehistory />
        </>}
    </div>
}

export const getServerSideProps = initializeStore;