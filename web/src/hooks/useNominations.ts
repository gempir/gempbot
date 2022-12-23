import { useEffect, useRef, useState } from "react";
import { doFetch, Method } from "../service/doFetch";
import { useStore } from "../store";

interface NominationVote {
    EmoteID: string
    ChannelTwitchID: string
    VotedBy: string
}

interface RawNomination {
    EmoteID: string
    ChannelTwitchID: string
    Votes: Array<NominationVote>
    EmoteCode: string
    NominatedBy: string
    CreatedAt: string
    UpdatedAt: string
}

export type Nomination = RawNomination & {
    CreatedAt: Date,
    UpdatedAt: Date,
}

interface Return {
    nominations: Array<Nomination>,
    fetch: () => void,
    makeVote: (emoteID: string) => void,
    loading: boolean,
}

export function useNominations(channel: string): Return {
    const [nominations, setBlocks] = useState<Array<Nomination>>([]);
    const [loading, setLoading] = useState(false);
    const apiBaseUrl = useStore(state => state.apiBaseUrl);
    const scToken = useStore(state => state.scToken);

    const fetchNominations = () => {
        setLoading(true);

        const endPoint = "/api/nominations";
        const searchParams = new URLSearchParams();
        searchParams.append("channel", channel);
        doFetch({ apiBaseUrl }, Method.GET, endPoint, searchParams).then((resp) => {
            return resp
        }).then(rawNoms => setBlocks(rawNoms.map((rawNom: RawNomination) => ({ ...rawNom, CreatedAt: new Date(rawNom.CreatedAt), UpdatedAt: new Date(rawNom.UpdatedAt) }))))
            .then(() => setLoading(false)).catch(err => {
                if (err.message !== "Page changed") {
                    throw err;
                }
            });
    };

    const makeVote = (emoteID: string) => {
        setLoading(true);

        const endPoint = "/api/nominations/vote";
        const searchParams = new URLSearchParams();
        searchParams.append("channel", channel);
        searchParams.append("emoteID", emoteID);
        doFetch({ apiBaseUrl, scToken }, Method.POST, endPoint, searchParams).then(() => setLoading(false)).catch(err => {}).finally(fetchNominations);
    };

    const interval = useRef<NodeJS.Timeout | null>(null);

    useEffect(() => {
        interval.current = setInterval(() => {
            fetchNominations();
        }, 10000);
        return () => {
            if (interval.current) {
                clearInterval(interval.current);
            }
        };
    }, []);

    useEffect(fetchNominations, []);

    return {
        nominations: nominations,
        fetch: fetchNominations,
        makeVote: makeVote,
        loading: loading
    };
}