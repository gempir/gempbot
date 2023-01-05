import dayjs, { Dayjs } from "dayjs";
import { useEffect, useState } from "react";
import { doFetch, Method } from "../service/doFetch";
import { useStore } from "../store";

interface NominationVote {
    EmoteID: string
    ChannelTwitchID: string
    VoteBy: string
}

interface RawNomination {
    EmoteID: string
    ChannelTwitchID: string
    Votes: Array<NominationVote>
    Downvotes: Array<NominationVote>
    EmoteCode: string
    NominatedBy: string
    CreatedAt: string
    UpdatedAt: string
}

export type Nomination = RawNomination & {
    CreatedAt: Dayjs,
    UpdatedAt: Dayjs,
}

interface Return {
    nominations: Array<Nomination>,
    fetch: () => void,
    vote: (emoteID: string) => void,
    unvote: (emoteID: string) => void,
    block: (emoteID: string) => void,
    downvote: (emoteID: string) => void,
    undownvote: (emoteID: string) => void,
    loading: boolean,
}

export function useNominations(channel: string): Return {
    const [nominations, setBlocks] = useState<Array<Nomination>>([]);
    const [loading, setLoading] = useState(false);
    const apiBaseUrl = useStore(state => state.apiBaseUrl);
    const managing = useStore(state => state.managing);
    const scToken = useStore(state => state.scToken);

    const fetchNominations = () => {
        setLoading(true);

        const endPoint = "/api/nominations";
        const searchParams = new URLSearchParams();
        searchParams.append("channel", channel);
        doFetch({ apiBaseUrl }, Method.GET, endPoint, searchParams).then((resp) => {
            return resp
        }).then(rawNoms => setBlocks(
            rawNoms.map(
                (rawNom: RawNomination) => (
                    {
                        ...rawNom,
                        CreatedAt: dayjs(rawNom.CreatedAt),
                        UpdatedAt: dayjs(rawNom.UpdatedAt)
                    }
                )
            )))
            .then(() => setLoading(false)).catch(err => {
                if (err.message !== "Page changed") {
                    throw err;
                }
            });
    };

    const vote = (emoteID: string) => {
        setLoading(true);

        const endPoint = "/api/nominations/vote";
        const searchParams = new URLSearchParams();
        searchParams.append("channel", channel);
        searchParams.append("emoteID", emoteID);
        doFetch({ apiBaseUrl, scToken }, Method.POST, endPoint, searchParams).then(() => setLoading(false)).catch(err => { }).finally(fetchNominations);
    };

    const unvote = (emoteID: string) => {
        setLoading(true);

        const endPoint = "/api/nominations/vote";
        const searchParams = new URLSearchParams();
        searchParams.append("channel", channel);
        searchParams.append("emoteID", emoteID);
        doFetch({ apiBaseUrl, scToken }, Method.DELETE, endPoint, searchParams).then(() => setLoading(false)).catch(err => { }).finally(fetchNominations);
    };

    const downvote = (emoteID: string) => {
        setLoading(true);

        const endPoint = "/api/nominations/downvote";
        const searchParams = new URLSearchParams();
        searchParams.append("channel", channel);
        searchParams.append("emoteID", emoteID);
        doFetch({ apiBaseUrl, scToken }, Method.POST, endPoint, searchParams).then(() => setLoading(false)).catch(err => { }).finally(fetchNominations);
    };

    const undownvote = (emoteID: string) => {
        setLoading(true);

        const endPoint = "/api/nominations/downvote";
        const searchParams = new URLSearchParams();
        searchParams.append("channel", channel);
        searchParams.append("emoteID", emoteID);
        doFetch({ apiBaseUrl, scToken }, Method.DELETE, endPoint, searchParams).then(() => setLoading(false)).catch(err => { }).finally(fetchNominations);
    };

    const block = (emoteID: string) => {
        setLoading(true);

        const endPoint = "/api/nominations/block";
        const searchParams = new URLSearchParams();
        searchParams.append("channel", channel);
        searchParams.append("emoteID", emoteID);
        doFetch({ apiBaseUrl, scToken, managing }, Method.DELETE, endPoint, searchParams).then(() => setLoading(false)).catch(err => { }).finally(fetchNominations);
    };

    useEffect(fetchNominations, []);

    return {
        nominations: nominations,
        fetch: fetchNominations,
        vote: vote,
        unvote: unvote,
        block: block,
        downvote: downvote,
        undownvote: undownvote,
        loading: loading
    };
}