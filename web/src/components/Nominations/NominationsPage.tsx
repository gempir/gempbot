import { useElection } from "../../hooks/useElection";
import { NominationsView } from "./NominationsView";

export function NominationsPage({ channel }: { channel: string }): JSX.Element {
    const [election] = useElection(channel);

    return <div className="p-4 w-full"><NominationsView channel={channel} election={election} /></div>;
}