import { useUserConfig } from "../../hooks/useUserConfig";
import { BttvForm } from "./RewardForms/BttvForm";
import { SevenTvForm } from "./RewardForms/SevenTvForm";

export function Rewards() {
    const [userCfg] = useUserConfig();
    if (!userCfg) {
        return null;
    }

    return <div className="p-4 flex gap-5 items-start">
        <BttvForm userConfig={userCfg} />
        <SevenTvForm userConfig={userCfg} />
    </div>;
}