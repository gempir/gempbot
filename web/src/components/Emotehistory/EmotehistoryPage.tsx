import { useParams } from "react-router-dom";
import { Emotehistory } from "../Home/Emotehistory";

export function EmotehistoryPage() {
    const { channel } = useParams<{channel: string}>();

    return <div className="p-4">
        <Emotehistory channel={channel} />
    </div>;
}