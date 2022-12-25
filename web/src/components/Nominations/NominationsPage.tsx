import { NominationsView } from "./NominationsView";

export function NominationsPage({ channel }: { channel: string }): JSX.Element {
    return <div className="p-4 w-full"><NominationsView channel={channel} /></div>;
}