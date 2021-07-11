import { useTitle } from "react-use";
import { useUserConfig } from "../../hooks/useUserConfig";
import { BttvForm } from "./RewardForms/BttvForm";

export function Rewards() {
    useTitle("bitraft - Rewards");

    const [userCfg] = useUserConfig();
    if (!userCfg) {
        return null;
    }

    return <div className="p-4">
        <BttvForm userConfig={userCfg} />
    </div>;
}