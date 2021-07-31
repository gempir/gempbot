import { Link } from "react-router-dom";
import { useUserConfig } from "../../hooks/useUserConfig";
import { Gift } from "../../icons/Gift";
import { House } from "../../icons/House";
import { Settings } from "../../icons/Settings";
import { BotManager } from "./BotManager";
import { Login } from "./Login";
import { Managing } from "./Managing";

export function Sidebar() {
    const [userConfig, setUserConfig] = useUserConfig();

    return <div className="p-4 bg-gray-800 px-6 shadow flex flex-col relative h-screen">
        <Login />
        <BotManager userConfig={userConfig!} setUserConfig={setUserConfig} />
        <Managing userConfig={userConfig} />
        <Link to="/" className="flex gap-2 items-center py-4 justify-start hover:text-blue-500">
            <House /> Home
        </Link>
        <Link to="/rewards" className="flex gap-2 items-center py-4 justify-start hover:text-blue-500">
            <Gift /> Rewards
        </Link>
        <Link to="/permissions" className="flex gap-2 items-center py-4 justify-start hover:text-blue-500">
            <Settings /> Permissions
        </Link>
        <Link to="/privacy" className="absolute bottom-3 text-center left-0 right-0 mx-auto hover:text-blue-500">
            Privacy
        </Link>
    </div>;
}