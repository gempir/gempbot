import { AdjustmentsHorizontalIcon, ChatBubbleLeftIcon, EyeIcon, GiftIcon, HomeIcon, NoSymbolIcon, PhotoIcon } from "@heroicons/react/24/solid";
import Link from "next/link";
import { useUserConfig } from "../../hooks/useUserConfig";
import { useStore } from "../../store";
import { Login } from "./Login";
import { Managing } from "./Managing";
import { NavLink } from "@mantine/core";
import { usePathname } from "next/navigation";

export function Sidebar() {
    const [userConfig] = useUserConfig();
    const loggedIn = useStore(state => Boolean(state.scToken));
    const pathname = usePathname();

    if (!loggedIn) {
        return null;
    }

    return (
        <div className="py-4 bg-gray-800 shadow flex flex-col relative h-screen items-center">
            <Login />
            {loggedIn && <>
                <Managing userConfig={userConfig} />
                <NavLink
                    href="/"
                    active={pathname === "/"}
                    label="Home"
                    leftSection={<HomeIcon className="h-6" />}
                    component={Link}
                />
                <NavLink
                    href="/bot"
                    active={pathname.startsWith("/bot")}
                    label="Bot"
                    leftSection={<ChatBubbleLeftIcon className="h-6" />}
                    component={Link}
                />
                <NavLink
                    href="/overlay"
                    active={pathname.startsWith("/overlay")}
                    label="Overlays"
                    leftSection={<PhotoIcon className="h-6" />}
                    component={Link}
                />
                <NavLink
                    href="/rewards"
                    active={pathname.startsWith("/rewards")}
                    label="Rewards"
                    leftSection={<GiftIcon className="h-6" />}
                    component={Link}
                />
                <NavLink
                    href="/permissions"
                    active={pathname.startsWith("/permissions")}
                    label="Permissions"
                    leftSection={<AdjustmentsHorizontalIcon className="h-6" />}
                    component={Link}
                />
                <NavLink
                    href="/blocks"
                    active={pathname.startsWith("/blocks")}
                    label="Blocks"
                    leftSection={<NoSymbolIcon className="h-6" />}
                    component={Link}
                />
            </>}
        </div>
    );
}