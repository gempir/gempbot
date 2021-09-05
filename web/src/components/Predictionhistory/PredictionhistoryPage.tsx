import { useParams } from "react-router-dom";
import { PredictionLog } from "../Home/PredictionLog";

export function PredictionhistoryPage() {
    const { channel } = useParams<{channel: string}>();

    return <div className="p-4 w-full"> 
        <PredictionLog channel={channel} />
    </div>;
}