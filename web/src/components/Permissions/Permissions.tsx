import { useUserConfig } from "../../hooks/useUserConfig";
import { UserPermissions } from "./UserPermissions";

export function Permissions() {
    const [userCfg, setUserConfig] = useUserConfig();
    if (!userCfg) {
        return null;
    }

    return <div className="p-4">
        <UserPermissions userConfig={userCfg} setUserConfig={setUserConfig} />
    </div>;
}