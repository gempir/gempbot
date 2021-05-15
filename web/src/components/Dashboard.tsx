import { useUserConfig } from "../hooks/useUserConfig";
import { Menu } from "./Menu";

export function Dashboard() {
    const [userCfg, setUserConfig] = useUserConfig();

    return <div>
        {userCfg && <Menu userConfig={userCfg} setUserConfig={setUserConfig} />}
        {/* {userCfg && <>
            <BttvForm />
        </>} */}
    </div>
}