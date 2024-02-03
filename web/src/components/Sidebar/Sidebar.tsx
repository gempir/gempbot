import { AdjustmentsHorizontalIcon, ChatBubbleLeftIcon, GiftIcon, HomeIcon, NoSymbolIcon, PhotoIcon } from "@heroicons/react/24/solid";
import Link from "next/link";
import { useUserConfig } from "../../hooks/useUserConfig";
import { useStore } from "../../store";
import { Login } from "./Login";
import { Managing } from "./Managing";

export function Sidebar() {
    const [userConfig] = useUserConfig();
    const isDev = useStore(state => state.baseUrl).includes("localhost");
    const loggedIn = useStore(state => Boolean(state.scToken));

    if (!loggedIn) {
        return null;
    }

    return <div className="p-4 bg-gray-800 px-6 shadow flex flex-col relative h-screen">
        <Login />
        {loggedIn && <>
            <Managing userConfig={userConfig} />
            <Link href="/">
                <a className="flex gap-2 items-center py-4 justify-start hover:text-blue-500 "><HomeIcon className="h-6" /> Home</a>
            </Link>
            <Link href="/bot" >
                <a className="flex gap-2 items-center py-4 justify-start hover:text-blue-500"><ChatBubbleLeftIcon className="h-6" /> Bot</a>
            </Link>
            {isDev && <Link href="/overlays" >
                <a className="flex gap-2 items-center py-4 justify-start hover:text-blue-500"><PhotoIcon className="h-6" /> Overlays</a>
            </Link>}
            <Link href="/rewards" >
                <a className="flex gap-2 items-center py-4 justify-start hover:text-blue-500"><GiftIcon className="h-6" /> Rewards</a>
            </Link>
            <Link href="/permissions" >
                <a className="flex gap-2 items-center py-4 justify-start hover:text-blue-500"><AdjustmentsHorizontalIcon className="h-6" /> Permissions</a>
            </Link>
            <Link href="/blocks" >
                <a className="flex gap-2 items-center py-4 justify-start hover:text-blue-500"><NoSymbolIcon className="h-6" /> Blocks</a>
            </Link>
        </>}
    </div>;
}