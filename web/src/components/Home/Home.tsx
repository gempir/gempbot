import { useStore } from "../../store";
import { EmoteLogPage } from "../EmoteLog/EmoteLogPage";

export function Home() {
    const scTokenContent = useStore(state => state.scTokenContent);
    const managing = useStore(state => state.managing);
    const channel = managing ?? scTokenContent?.Login;

    if (!channel) {
        return null;
    }

    return <div>
        <EmoteLogPage channel={channel} />
    </div>
}