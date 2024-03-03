import { AdjustmentsHorizontalIcon, ChatBubbleLeftIcon, EyeIcon, GiftIcon, HomeIcon, NoSymbolIcon, PhotoIcon } from "@heroicons/react/24/solid";
import Link from "next/link";
import { useUserConfig } from "../../hooks/useUserConfig";
import { useStore } from "../../store";
import { Login } from "./Login";
import { Managing } from "./Managing";

export function Sidebar() {
    const [userConfig] = useUserConfig();
    const loggedIn = useStore(state => Boolean(state.scToken));

    if (!loggedIn) {
        return null;
    }

    return (
        <div className="p-4 bg-gray-800 px-6 shadow flex flex-col relative h-screen">
            <Login />
            {loggedIn && <>
                <Managing userConfig={userConfig} />
                <Link
                    href="/"
                    className="flex gap-2 items-center py-4 justify-start hover:text-blue-500 ">
                    <HomeIcon className="h-6" />Home
                </Link>
                <Link
                    href="/george"
                    className="flex gap-2 items-center py-4 justify-start hover:text-blue-500">
                    <EyeIcon className="h-6" />George
                </Link>
                <Link
                    href="/bot"
                    className="flex gap-2 items-center py-4 justify-start hover:text-blue-500">
                    <ChatBubbleLeftIcon className="h-6" />Bot
                </Link>
                <Link
                    href="/overlay"
                    className="flex gap-2 items-center py-4 justify-start hover:text-blue-500">
                    <PhotoIcon className="h-6" />Overlays
                </Link>
                <Link
                    href="/rewards"
                    className="flex gap-2 items-center py-4 justify-start hover:text-blue-500">
                    <GiftIcon className="h-6" />Rewards
                </Link>
                <Link
                    href="/permissions"
                    className="flex gap-2 items-center py-4 justify-start hover:text-blue-500">
                    <AdjustmentsHorizontalIcon className="h-6" />Permissions
                </Link>
                <Link
                    href="/blocks"
                    className="flex gap-2 items-center py-4 justify-start hover:text-blue-500">
                    <NoSymbolIcon className="h-6" />Blocks
                </Link>
            </>}
        </div>
    );
}