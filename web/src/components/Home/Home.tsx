import { useTitle } from "react-use";
import { EmoteHistory } from "./EmoteHistory";
import { PredictionLog } from "./PredictionLog";

export function Home() {
    useTitle("bitraft - Home");

    return <div className="p-4 w-full max-h-screen flex gap-4">
        <EmoteHistory />
        <PredictionLog />
    </div>
}