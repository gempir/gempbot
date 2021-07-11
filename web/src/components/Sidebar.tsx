import { Link } from "react-router-dom";
import { Settings } from "../icons/Settings";
import { Gift } from "../icons/Gift";
import { Login } from "./Login";
import { House } from "../icons/House";

export function Sidebar() {
    return <div className="p-4 bg-gray-800 w-48 rounded shadow flex flex-col relative h-screen">
        <Login />
        <div className="h-10" />
        <Link to="/" className="flex gap-2 items-center py-4 justify-start hover:text-gray-400">
            <House /> Home
        </Link>
        <Link to="/rewards" className="flex gap-2 items-center py-4 justify-start hover:text-gray-400">
            <Gift /> Rewards
        </Link>
        <Link to="/permissions" className="flex gap-2 items-center py-4 justify-start hover:text-gray-400">
            <Settings /> Permissions
        </Link>
        <Link to="/privacy" className="absolute bottom-3 text-center left-0 right-0 mx-auto hover:text-gray-400">
            Privacy
        </Link>
    </div>;
}