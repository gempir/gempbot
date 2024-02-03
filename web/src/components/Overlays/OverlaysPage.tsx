import { useUserConfig } from "../../hooks/useUserConfig";

export function OverlaysPage() {
    const [userCfg, setUserConfig, , loading, errorMessage] = useUserConfig();
    if (!userCfg) {
        return null;
    }

    return <div className="p-4">
        Overlays Here
    </div>;
}