import { useUserConfig } from "../../hooks/useUserConfig";
import { useStore } from "../../store";
import { BotManager } from "./BotManager";
import { Login } from "./Login";
import { Managing } from "./Managing";
import Link from "next/link";
import React from "react";
import { House } from "../../icons/House";
import { Gift } from "../../icons/Gift";
import { Settings } from "../../icons/Settings";

export function Sidebar() {
    const [userConfig, setUserConfig] = useUserConfig();
    const loggedIn = useStore(state => Boolean(state.scToken));

    return <div className="p-4 bg-gray-800 px-6 shadow flex flex-col relative h-screen">
        <Login />
        {loggedIn && <>
            <BotManager userConfig={userConfig!} setUserConfig={setUserConfig} />
            <Managing userConfig={userConfig} />
            <Link href="/">
                <a className="flex gap-2 items-center py-4 justify-start hover:text-blue-500 "><House /> Home</a>
            </Link>
            <Link href="/rewards" >
                <a className="flex gap-2 items-center py-4 justify-start hover:text-blue-500"><Gift /> Rewards</a>
            </Link>
            <Link href="/permissions" >
                <a className="flex gap-2 items-center py-4 justify-start hover:text-blue-500"><Settings /> Permissions</a>
            </Link>
            <div className="absolute bottom-3 text-center left-0 right-0 mx-auto hover:text-blue-500">
                <Link href="/privacy">
                    <a>Privacy</a>
                </Link>
            </div>
        </>}
    </div>;
}