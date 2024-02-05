import dynamic from "next/dynamic";
import { useUserConfig } from "../../hooks/useUserConfig";
const Editor = dynamic(async () => (await import('./../Overlay/Editor')).Editor, { ssr: false })

export function OverlaysPage() {
    const [userCfg, setUserConfig, , loading, errorMessage] = useUserConfig();
    if (!userCfg) {
        return null;
    }

    return <div className="relative w-full h-[100vh]">
        <Editor />
    </div>;
}