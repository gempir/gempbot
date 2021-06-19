import { useTitle } from "react-use";
import { useUserConfig } from "../../hooks/useUserConfig";
import { Menu } from "./Menu";
import { BttvForm } from "./RewardForms/BttvForm";

export function Dashboard() {
    useTitle("bitraft - Dashboard");
    const [userCfg, setUserConfig] = useUserConfig();
    if (!userCfg) {
        return null;
    }

    return <div>
        <Menu userConfig={userCfg} setUserConfig={setUserConfig} />
        <BttvForm userConfig={userCfg} />
    </div>
}