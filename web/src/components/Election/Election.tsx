import { useElection } from "../../hooks/useElection";
import { useUserConfig } from "../../hooks/useUserConfig";
import { useStore } from "../../store";
import { NominationsView } from "../Nominations/NominationsView";
import { ElectionForm } from "./ElectionForm";

export function Election() {
    const scTokenContent = useStore(state => state.scTokenContent);
    const managing = useStore(state => state.managing);
    const channel = managing ?? scTokenContent?.Login;
    const [election, setElection, deleteElection, electionErrorMessage, electionLoading] = useElection();

    const [userCfg, setUserConfig, , loading, errorMessage] = useUserConfig();
    if (!userCfg) {
        return null;
    }

    return <div className="p-4 flex gap-3">
        <div>
            <ElectionForm election={election} setElection={setElection} deleteElection={deleteElection} electionErrorMessage={electionErrorMessage} electionLoading={electionLoading} />
        </div>
        {channel && <NominationsView showVotes channel={channel} election={election} />}
    </div>;
}