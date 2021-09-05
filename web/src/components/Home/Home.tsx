import { useTitle } from "react-use";
import { Emotehistory } from "./Emotehistory";
import { PredictionLog } from "./PredictionLog";

export function Home() {
    useTitle("bitraft - Home");

    return <div className="p-4 w-full max-h-screen flex gap-4">
        <Emotehistory />
        <PredictionLog />
    </div>
}