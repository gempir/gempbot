import { useRef } from "react";
import { useEffect, useState } from "react";
import { PredictionLog } from "../model/PredictionLog";
import { doFetch, Method } from "../service/doFetch";
import { store } from "../store";

export interface Outcome {
    ID: string;
    PredictionID: string;
    Title: string;
    Color: string;
    Users: number;
    ChannelPoints: number;
}

export interface RawPredictionLog {
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

const PAGE_SIZE = 20;

export function usePredictionLogs(channel?: string): [Array<PredictionLog>, () => void, boolean, number, () => void, () => void] {
    const [page, setPage] = useState(1);
    const pageRef = useRef(page);
    pageRef.current = page;

    const [logs, setLogs] = useState<Array<PredictionLog>>([]);
    const [loading, setLoading] = useState(false);
    const managing = store.useState(s => s.managing);

    const fetchPredictions = () => {
        setLoading(true);

        const currentPage = pageRef.current;
        
        const endPoint = "/api/predictionhistory";
        const searchParams = new URLSearchParams();
        if (channel) {
            searchParams.append("channel", channel);
        }
        searchParams.append("page", page.toString());
        doFetch(Method.GET, endPoint, searchParams).then((resp) => {
            if (currentPage !== pageRef.current) {
                throw new Error("Page changed");
            }

            return resp
        }).then((logs: Array<RawPredictionLog>) => setLogs(logs.map((log: RawPredictionLog) => PredictionLog.fromObject(log)))).then(() => setLoading(false)).catch(err => {
            if (err.message !== "Page changed") {
                throw err;
            }
        });
    };

    // eslint-disable-next-line react-hooks/exhaustive-deps
    useEffect(fetchPredictions, [managing, page]);

    return [logs, fetchPredictions, loading, page, () => logs.length === PAGE_SIZE ? setPage(page + 1) : undefined, () => page > 1 ? setPage(page - 1) : undefined];
}