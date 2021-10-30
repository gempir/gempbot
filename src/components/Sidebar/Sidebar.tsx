import { AdjustmentsIcon, BanIcon, GiftIcon, HomeIcon } from "@heroicons/react/solid";
import Link from "next/link";
import React from "react";
import { useUserConfig } from "../../hooks/useUserConfig";
import { useStore } from "../../store";
import { BotManager } from "./BotManager";
import { EventSubManager } from "./EventSubManager";
import { Login } from "./Login";
import { Managing } from "./Managing";

export function Sidebar() {
    const [userConfig, setUserConfig, , loading] = useUserConfig();
    const loggedIn = useStore(state => Boolean(state.scToken));

    return <div className="p-4 bg-gray-800 px-6 shadow flex flex-col relative h-screen">
        <Login />
        {loggedIn && <>
            <BotManager userConfig={userConfig!} setUserConfig={setUserConfig} userConfigLoading={loading} />
            <EventSubManager />
            <Managing userConfig={userConfig} />
            <Link href="/">
                <a className="flex gap-2 items-center py-4 justify-start hover:text-blue-500 "><HomeIcon className="h-6" /> Home</a>
            </Link>
            <Link href="/rewards" >
                <a className="flex gap-2 items-center py-4 justify-start hover:text-blue-500"><GiftIcon className="h-6" /> Rewards</a>
            </Link>
            <Link href="/permissions" >
                <a className="flex gap-2 items-center py-4 justify-start hover:text-blue-500"><AdjustmentsIcon className="h-6" /> Permissions</a>
            </Link>
            <Link href="/blocks" >
                <a className="flex gap-2 items-center py-4 justify-start hover:text-blue-500"><BanIcon className="h-6" /> Blocks</a>
            </Link>
            <div className="absolute bottom-3 text-center left-0 right-0 mx-auto hover:text-blue-500">
                <Link href="/privacy">
                    <a>Privacy</a>
                </Link>
            </div>
        </>}
    </div>;
}