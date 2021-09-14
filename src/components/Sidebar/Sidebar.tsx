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
    const scToken = useStore(state => state.scToken);


    return <div className="p-4 bg-gray-800 px-6 shadow flex flex-col relative h-screen">
        <Login />
        {scToken && <>
            <BotManager userConfig={userConfig!} setUserConfig={setUserConfig} />
            <Managing userConfig={userConfig} />
            <div className="flex gap-2 items-center py-4 justify-start hover:text-blue-500">
                <Link href="/">
                    <House /> Home
                </Link>
            </div>
            <div className="flex gap-2 items-center py-4 justify-start hover:text-blue-500">
                <Link href="/rewards" >
                    <Gift /> Rewards
                </Link>
            </div>
            <div className="flex gap-2 items-center py-4 justify-start hover:text-blue-500">
                <Link href="/permissions" >
                    <Settings /> Permissions
                </Link>
            </div>
            <div className="absolute bottom-3 text-center left-0 right-0 mx-auto hover:text-blue-500">
                <Link href="/privacy">
                    Privacy
                </Link>
            </div>
        </>}
    </div>;
}