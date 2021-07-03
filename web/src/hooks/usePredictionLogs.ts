import { useEffect, useState } from "react";
import { doFetch, Method } from "../service/doFetch";
import { store } from "../store";

export interface PredictionLog {
    ID: string;
    OwnerTwitchID: string;
    Title: string;
    WinningOutcomeID: string;
    Status: string;
    StartedAt: Date;
    LockedAt: Date;
    EndedAt: Date;
    Outcomes: Outcome[];
}

export interface Outcome {
    ID: string;
    PredictionID: string;
    Title: string;
    Color: string;
    Users: number;
    ChannelPoints: number;
}

interface RawPredictionLog {
    ID: string;
    OwnerTwitchID: string;
    Title: string;
    WinningOutcomeID: string;
    Status: string;
    StartedAt: string;
    LockedAt: string;
    EndedAt: string;
    Outcomes: Outcome[];
}

export function usePredictionLogs(): [Array<PredictionLog>, () => void] {
    const [logs, setLogs] = useState<Array<PredictionLog>>([]);
    const managing = store.useState(s => s.managing);

    const fetchConfig = () => {
        let endPoint = "/api/prediction";
        if (managing) {
            endPoint += `?managing=${managing}`;
        }
        doFetch(Method.GET, endPoint).then((logs) => setLogs(logs.map((log: RawPredictionLog) => ({
            ...log,
            StartedAt: new Date(log.StartedAt),
            LockedAt: new Date(log.LockedAt),
            EndedAt: new Date(log.EndedAt)
        }))))
    };

    // eslint-disable-next-line react-hooks/exhaustive-deps
    useEffect(fetchConfig, [managing]);

    return [logs, fetchConfig];
}