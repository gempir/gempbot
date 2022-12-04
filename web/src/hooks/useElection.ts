import { useEffect, useState } from "react";
import { doFetch, Method, RejectReason } from "../service/doFetch";
import { useStore } from "../store";
import { Election } from "../types/Election";

const defaultElection: Election = {
    Hours: 24,
    NominationCost: 1000,
}

export function useElection(): [Election, (election: Election) => void, () => void, string | null, boolean] {
    const [election, setElection] = useState<Election>(defaultElection);
    const [errorMessage, setErrorMessage] = useState<string | null>(null);
    const [loading, setLoading] = useState(false);
    const managing = useStore(state => state.managing);
    const apiBaseUrl = useStore(state => state.apiBaseUrl);
    const scToken = useStore(state => state.scToken);

    const fetchElection = () => {
        setLoading(true);
        doFetch({ apiBaseUrl, managing, scToken }, Method.GET, "/api/election").then(resp => setElection(
            {
                ...resp,
                CreatedAt: new Date(resp.CreatedAt ?? null),
                UpdatedAt: new Date(resp.UpdatedAt ?? null),
                LastRunAt: new Date(resp.LastRunAt ?? null),
            }
        )).catch(reason => {
            if (reason !== RejectReason.NotFound) {
                throw new Error(reason);
            }
        }).finally(() => setLoading(false));
    }

    useEffect(fetchElection, [managing]);

    const updateElection = (election: Election) => {
        setLoading(true);
        doFetch({ apiBaseUrl, managing, scToken }, Method.POST, "/api/election", undefined, election).then(() => setErrorMessage(null)).then(fetchElection).catch(setErrorMessage).finally(() => setLoading(false));
    }

    const deleteElection = () => {
        setLoading(true);
        doFetch({ apiBaseUrl, managing, scToken }, Method.DELETE, "/api/election").then(() => setErrorMessage(null)).then(fetchElection).catch(setErrorMessage).finally(() => setLoading(false));
    }

    return [election, updateElection, deleteElection, errorMessage, loading];
}
