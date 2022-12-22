import { useUserConfig } from "../../hooks/useUserConfig";
import { ElectionForm } from "./ElectionForm";

export function Election() {
    const [userCfg, setUserConfig, , loading, errorMessage] = useUserConfig();
    if (!userCfg) {
        return null;
    }

    return <div className="p-4">
        <ElectionForm />
    </div>;
}