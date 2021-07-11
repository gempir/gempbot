import { useTitle } from "react-use";
import { useUserConfig } from "../../hooks/useUserConfig";
import { Menu } from "./Menu";
import { PredictionLog } from "./PredictionLog";

export function Dashboard() {
    useTitle("bitraft - Dashboard");
    const [userCfg, setUserConfig] = useUserConfig();
    if (!userCfg) {
        return null;
    }

    return <div className="p-4 w-full max-h-screen overflow-y-scroll">
        <Menu userConfig={userCfg} setUserConfig={setUserConfig} />
        <div className="flex w-full">
            <PredictionLog />
        </div>
    </div>
}