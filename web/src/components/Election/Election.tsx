import { useUserConfig } from "../../hooks/useUserConfig";
import { useStore } from "../../store";
import { NominationsView } from "../Nominations/NominationsView";
import { ElectionForm } from "./ElectionForm";

export function Election() {
    const scTokenContent = useStore(state => state.scTokenContent);
    const [userCfg, setUserConfig, , loading, errorMessage] = useUserConfig();
    if (!userCfg) {
        return null;
    }


    return <div className="p-4 flex gap-3">
        <div>
            <ElectionForm />
        </div>
        {scTokenContent?.Login && <NominationsView channel={scTokenContent?.Login} />}
    </div>;
}