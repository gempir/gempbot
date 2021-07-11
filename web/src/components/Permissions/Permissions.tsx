import { useTitle } from "react-use";
import { useUserConfig } from "../../hooks/useUserConfig";
import { EditorManager } from "./EditorManager";
import { UserPermissions } from "./UserPermissions";

export function Permissions() {
    useTitle("bitraft - Permissions");

    const [userCfg, setUserConfig] = useUserConfig();
    if (!userCfg) {
        return null;
    }

    return <div className="p-4">
        <EditorManager userConfig={userCfg} setUserConfig={setUserConfig} />
        <UserPermissions userConfig={userCfg} setUserConfig={setUserConfig} />
    </div>;
}