import { useUserConfig } from "../../hooks/useUserConfig";
import { NewElection } from "./NewElection";

export function Election() {
    const [userCfg, setUserConfig, , loading, errorMessage] = useUserConfig();
    if (!userCfg) {
        return null;
    }

    return <div className="p-4">
        <NewElection />
    </div>;
}