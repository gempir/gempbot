import { useTitle } from "react-use";
import { PredictionLog } from "./PredictionLog";

export function Home() {
    useTitle("bitraft - Home");

    return <div className="p-4 w-full max-h-screen overflow-y-scroll">
        <PredictionLog />
    </div>
}