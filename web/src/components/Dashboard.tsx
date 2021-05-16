import { useState } from "react";
import { useUserConfig } from "../hooks/useUserConfig";
import { Menu } from "./Menu";
import { BttvForm } from "./RewardForms/BttvForm";

export function Dashboard() {
    const [renderKey, setRenderKey] = useState(1);
    const [userCfg, setUserConfig] = useUserConfig(undefined, undefined, () => setRenderKey(renderKey + 1));

    if (!userCfg) {
        return null;
    }

    // force re-mount when switching the channel to manage, to re-render forms and their defaultValues
    return <div key={renderKey}>
        <Menu userConfig={userCfg} setUserConfig={setUserConfig} />
        <BttvForm userConfig={userCfg} setUserConfig={setUserConfig} />
    </div>
}